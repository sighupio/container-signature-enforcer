package notary

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/sighupio/opa-notary-connector/config"
	"github.com/theupdateframework/notary"
	"github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/trustpinning"
	"github.com/theupdateframework/notary/tuf/data"
)

func NewFileCachedRepository(c *config.GlobalConfig, repo *config.Repository, ref *Reference, log *logrus.Entry) (*client.Repository, error) {
	contextLogger := log.WithFields(logrus.Fields{"image": ref.original, "server": repo.Trust.TrustServer})
	contextLogger.WithField("signers", repo.Trust.Signers).Debug("Checking image against server for signers")
	// initialize the repo
	r, err := client.NewFileCachedRepository(
		c.TrustRootDir,
		data.GUN(ref.name),
		repo.Trust.TrustServer,
		makeHubTransport(repo.Trust.TrustServer, ref.name, log),
		nil, //no need for passRetriever ATM
		//TODO: pass the notary CA explicitly via conf
		trustpinning.TrustPinConfig{},
	)
	if err != nil {
		contextLogger.WithError(err).Error("Error creating repository")
	}
	return &r, err
}

// returns the sha of an image in a given trust server
func CheckImage(ref *Reference, rootDir string, repo *config.Repository, client *client.Repository, log *logrus.Entry) (string, error) {
	contextLogger := log.WithFields(logrus.Fields{"image": ref, "server": repo.Trust.TrustServer})

	// build the roles from the signers
	rolelist := []data.RoleName{}
	rolesToPublicKeys := map[data.RoleName]data.PublicKey{}
	//TODO: do it once per config
	rolesFound := map[data.RoleName]bool{}
	for _, signer := range repo.Trust.Signers {
		role := data.RoleName(signer.Role)
		rolelist = append(rolelist, role)
		keyFromConfig, err := signer.GetPEM(contextLogger)

		if err != nil || keyFromConfig == nil {
			contextLogger.WithField("signer", signer).WithError(err).Debug("Error parsing public key")
			return "", err
		}
		contextLogger.WithFields(logrus.Fields{"signer": signer, "parsedPublicKey": keyFromConfig}).Debug("returning parsed public key")

		rolesToPublicKeys[role] = *keyFromConfig
		rolesFound[role] = false
	}

	/////////////////////////////////////// modified from Portieris codebase
	targets, err := (*client).GetAllTargetMetadataByName(ref.tag)

	contextLogger.WithFields(logrus.Fields{"ref": ref, "targets": targets}).Debug("Retrieved targets for image from server")
	if err != nil {
		contextLogger.WithError(err).Error("GetAllTargetMetadataByName returned an error")
		return "", err
	}

	if len(targets) == 0 {
		contextLogger.Error("No signed targets found")
		return "", fmt.Errorf("No signed targets found")
	}

	var digest []byte // holds digest of the signed image

	if len(rolelist) == 0 {
		// if no signer specified, no way to decide between the available targets, accept the last one
		for _, target := range targets {
			digest = target.Target.Hashes[notary.SHA256]
		}
		contextLogger.WithField("digest", digest).Debug("RoleList length == 0, returning digest", digest)
	} else {
		contextLogger.WithFields(logrus.Fields{"rolelist": rolelist, "targets": targets}).Debug("Looking for roles iterating over targets")
		// filter out targets signed by not required roles
		for _, target := range targets { // iterate over each target

			// See if a signer was specified for this target
			if keyFromConfig, ok := rolesToPublicKeys[target.Role.Name]; ok {
				if keyFromConfig != nil {
					// Assuming public key is in PEM format and not encoded any further
					contextLogger = contextLogger.WithField("role", target.Role.Name)
					contextLogger.WithFields(logrus.Fields{"keyID": keyFromConfig.ID(), "keys": target.Role.BaseRole.Keys}).Debug("Looking for key ID in keys")
					if _, ok := target.Role.BaseRole.Keys[keyFromConfig.ID()]; !ok {
						contextLogger.WithFields(logrus.Fields{"keyID": keyFromConfig.ID(), "keys": target.Role.BaseRole.ListKeyIDs()}).Error("KeyID not found in role key list")
						return "", fmt.Errorf("Public keys are different")
					}
					// We found a matching KeyID, so mark the role found in the map.
					contextLogger.WithField("keyID", keyFromConfig.ID()).Debug("found role with keyID")
					// store the digest of the latest signed release
					rolesFound[target.Role.Name] = true
				} else {
					contextLogger.WithField("role", target.Role.Name).Error("PublicKey not specified for role")
					return "", fmt.Errorf("PublicKey not specified for role %s", target.Role.Name)
				}

				// verify that the digest is consistent between all of the targets we care about
				if digest != nil && !bytes.Equal(digest, target.Target.Hashes[notary.SHA256]) {
					contextLogger.WithFields(logrus.Fields{"digest": digest, "target": target}).Error("Digest is different from that of target")
					return "", fmt.Errorf("Incompatible digest %s from that of target %+v", digest, target)
				} else {
					contextLogger.Debug("setting digest")
					digest = target.Target.Hashes[notary.SHA256]
				}
			}
		}

		//check all signatures from all specified roles have been found, overwise return error
		for role, found := range rolesFound {
			if !found {
				log.WithFields(logrus.Fields{"role": role, "key": rolesToPublicKeys[role]}).Error("Role not found with key")
				return "", fmt.Errorf("%s role not found with key %s", role, rolesToPublicKeys[role])
			}
		}
	}
	//////////////////////////////////////////////

	stringDigest := hex.EncodeToString(digest)
	contextLogger.WithField("digest", stringDigest).Debug("Returning digest for image")
	return stringDigest, nil
}

func makeHubTransport(server, image string, log *logrus.Entry) http.RoundTripper {
	base := http.DefaultTransport
	modifiers := []transport.RequestModifier{
		transport.NewHeaderRequestModifier(http.Header{
			"User-Agent": []string{"notary-admission-webhook"},
		}),
	}

	authTransport := transport.NewTransport(base, modifiers...)
	pingClient := &http.Client{
		Transport: authTransport,
		Timeout:   5 * time.Second,
	}
	req, err := http.NewRequest("GET", server+"/v2/", nil)
	if err != nil {
		log.WithError(err).WithField("server", server).Error("Error reading from notary server")
		return nil
	}

	challengeManager := challenge.NewSimpleManager()
	resp, err := pingClient.Do(req)
	if err != nil {
		log.WithError(err).WithField("server", server).Error("Error reading from notary server")
		return nil
	}

	defer resp.Body.Close()

	if err := challengeManager.AddResponse(resp); err != nil {
		log.WithError(err).WithField("server", server).Error("Error reading from notary server")
		return nil
	}

	tokenHandler := auth.NewTokenHandler(base, nil, image, "pull")
	modifiers = append(modifiers, auth.NewAuthorizer(challengeManager, tokenHandler, auth.NewBasicHandler(nil)))

	return transport.NewTransport(base, modifiers...)
}
