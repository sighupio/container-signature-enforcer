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
	"github.com/sighupio/opa-notary-connector/reference"
	"github.com/theupdateframework/notary"
	"github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/trustpinning"
	"github.com/theupdateframework/notary/tuf/data"
)

// AllTargetMetadataByNameGetter abstracts the only function we use of Notary, getting metadata by name for a target
type AllTargetMetadataByNameGetter interface {
	GetAllTargetMetadataByName(tag string) ([]client.TargetSignedStruct, error)
}

type Repository struct {
	rolesFound        map[data.RoleName]bool
	rolesToPublicKeys map[data.RoleName]data.PublicKey
	clientRepository  *AllTargetMetadataByNameGetter
	trustRootDir      string
	configRepository  *config.Repository
	log               *logrus.Entry
	reference         *reference.Reference
}

// NewWithGetter creates a Repository
func NewWithGetter(ref *reference.Reference, repo *config.Repository, getter *AllTargetMetadataByNameGetter, trustRootDir string, log *logrus.Entry) (*Repository, error) {
	no := Repository{
		configRepository:  repo,
		reference:         ref,
		log:               log,
		rolesFound:        make(map[data.RoleName]bool),
		rolesToPublicKeys: make(map[data.RoleName]data.PublicKey),
		clientRepository:  getter,
		trustRootDir:      trustRootDir,
	}

	return &no, nil
}

// New wraps NewWithGetter but then creates a FileCachedRepository as clientRepository, connecting to a real notary instance
func New(ref *reference.Reference, repo *config.Repository, trustRootDir string, log *logrus.Entry) (*Repository, error) {
	no, err := NewWithGetter(ref, repo, nil, trustRootDir, log)
	if err != nil {
		log.WithFields(logrus.Fields{
			"image":  ref,
			"server": repo.Trust.TrustServer,
		}).WithError(err).Error("failed image parsing")
	}
	err = no.newFileCachedRepository()
	if err != nil {
		log.WithFields(logrus.Fields{
			"image":  ref,
			"server": repo.Trust.TrustServer,
		}).WithError(err).Error("failed creating file cached repository")
		return nil, err
	}

	return no, nil
}

// getRolesFromSigners is an internal function to warm up the internal caches
func (no *Repository) getRolesFromSigners(signers []*config.Signer, log *logrus.Entry) (err error) {
	// build the roles from the signers
	for _, signer := range signers {
		role := data.RoleName(signer.Role)

		//TODO: parse multiple keys per same signer (list, not single key)
		keyFromConfig, err := signer.GetPEM(log)

		if err != nil || keyFromConfig == nil {
			log.WithField("signer", signer).WithError(err).Debug("Error parsing public key")
			return err
		}
		log.WithFields(logrus.Fields{"signer": signer, "parsedPublicKey": keyFromConfig}).Debug("returning parsed public key")

		no.rolesToPublicKeys[role] = *keyFromConfig
		no.rolesFound[role] = false
	}
	return nil
}

// returns the sha of an image in a given trust server
//ref *Reference, rootDir string, repo *config.Repository
func (no *Repository) GetSha() (string, error) {
	contextLogger := no.log.WithFields(logrus.Fields{"image": no.reference, "server": no.configRepository.Trust.TrustServer})

	err := no.getRolesFromSigners(no.configRepository.Trust.Signers, contextLogger)
	if err != nil {
		contextLogger.WithError(err).Error("getRolesFromSigners returned an error")
		return "", err
	}

	targets, err := (*no.clientRepository).GetAllTargetMetadataByName(no.reference.Tag)

	contextLogger.WithFields(logrus.Fields{"ref": no.reference, "targets": targets}).Debug("Retrieved targets for image from server")
	if err != nil {
		contextLogger.WithError(err).Error("GetAllTargetMetadataByName returned an error")
		return "", err
	}

	//target signers
	//0 0 => "", nil
	//0 m => "", error
	//n 0 => sha, nil
	//n m => check
	if len(targets) == 0 {
		if len(no.configRepository.Trust.Signers) == 0 {
			return "", nil
		} else {
			contextLogger.Error("No signed targets found")
			return "", fmt.Errorf("No signed targets found")
		}
	} else {
		var digest []byte // holds digest of the signed image
		if len(no.configRepository.Trust.Signers) == 0 {
			// if no signer specified, no way to decide between the available targets, accept the last one
			digest = targets[0].Target.Hashes[notary.SHA256]
			contextLogger.
				WithField("digest", digest).
				Debug("no.configRepository.Trust.Signers length == 0, returning digest")
		} else {
			// filter out targets signed by not required roles
			for _, target := range targets { // iterate over each target
				d, err := no.getShaFromTarget(&target, contextLogger)
				if err != nil {
					return "", err
				}
				if digest != nil && !bytes.Equal(digest, d) {
					contextLogger.
						WithFields(logrus.Fields{"digest": digest, "target": target}).
						Error("Digest is different from that of target")
					return "", fmt.Errorf("Incompatible digest from that of target")
				}
				digest = d
			}
			//check all signatures from all specified roles have been found, overwise return error
			for role, found := range no.rolesFound {
				if !found {
					contextLogger.
						WithFields(logrus.Fields{"role": role, "key": no.rolesToPublicKeys[role]}).
						Error("Role not found with key")
					return "", fmt.Errorf("Role not found with for a specified signer")
				}
			}
		}
		stringDigest := hex.EncodeToString(digest)
		contextLogger.WithField("digest", stringDigest).Debug("Returning digest for image")
		return stringDigest, nil
	}
}

func (no *Repository) getShaFromTarget(target *client.TargetSignedStruct, log *logrus.Entry) (digest []byte, err error) {
	log.WithFields(logrus.Fields{"signers": no.configRepository.Trust.Signers, "target": target}).Debug("Looking for roles iterating over targets")

	log = log.WithField("role", target.Role.Name)
	keyFromConfig, ok := no.rolesToPublicKeys[target.Role.Name]
	if !ok {
		return nil, nil
	}
	// Assuming public key is in PEM format and not encoded any further
	log.WithFields(logrus.Fields{"keyID": keyFromConfig.ID(), "keys": target.Role.BaseRole.Keys}).Debug("Looking for key ID in keys")
	if _, ok := target.Role.BaseRole.Keys[keyFromConfig.ID()]; !ok {
		log.WithFields(logrus.Fields{"keyID": keyFromConfig.ID(), "keys": target.Role.BaseRole.ListKeyIDs()}).Error("KeyID not found in role key list")
		return nil, fmt.Errorf("Public keys are different")
	}
	// We found a matching KeyID, so mark the role found in the map.
	log.WithField("keyID", keyFromConfig.ID()).Debug("found role with keyID")
	// store the digest of the latest signed release
	no.rolesFound[target.Role.Name] = true

	// verify that the digest is consistent between all of the targets we care about
	digest = target.Target.Hashes[notary.SHA256]
	log.WithField("sha256", digest).Debug("set digest")
	return digest, nil
}

// reference is notary lingo for image
func (no *Repository) newFileCachedRepository() error {
	contextLogger := no.log.WithFields(logrus.Fields{"image": no.reference.Original, "server": no.configRepository.Trust.TrustServer})
	contextLogger.WithField("signers", no.configRepository.Trust.Signers).Debug("Checking image against server for signers")
	// initialize the repo
	var r AllTargetMetadataByNameGetter
	r, err := client.NewFileCachedRepository(
		no.trustRootDir,
		data.GUN(no.reference.Name),
		no.configRepository.Trust.TrustServer,
		no.makeHubTransport(no.configRepository.Trust.TrustServer, no.reference.Name, contextLogger),
		nil, //no need for passRetriever ATM
		//TODO: pass the notary CA explicitly via conf
		trustpinning.TrustPinConfig{},
	)
	if err != nil {
		contextLogger.WithError(err).Error("Error creating repository")
	}
	no.clientRepository = &r
	return err
}

func (no *Repository) makeHubTransport(server, image string, log *logrus.Entry) http.RoundTripper {
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
