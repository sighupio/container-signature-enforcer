package config

import (
	"encoding/base64"
	"regexp"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/theupdateframework/notary/tuf/data"
	"github.com/theupdateframework/notary/tuf/utils"
)

type GlobalConfig struct {
	config *Config
	BindAddress,
	LogLevel,
	NotaryServer,
	TrustConfigPath,
	TrustRootDir string
	Mutex *sync.RWMutex
}

func NewGlobaConfig() *GlobalConfig {
	g := &GlobalConfig{}
	g.Mutex = new(sync.RWMutex)
	g.config = NewConfig()
	return g
}

func (g *GlobalConfig) SetConfig(c *Config) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	g.config = c
}

func (g *GlobalConfig) GetConfig() *Config {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	return g.config
}

// The main config structure containing informations regarding the image signature verification
type Config struct {
	Repositories       Repositories `mapstructure:"repositories"`
	sortedRepositories bool
	mutex              *sync.RWMutex
	repositoriesMutex  *sync.RWMutex
	//TODO: global signers that can be referenced by name from Trust, so that keys can be shared between repositories if needed
	//GlobalSigners []Signer
}

func NewConfig() *Config {
	return &Config{
		mutex:             new(sync.RWMutex),
		repositoriesMutex: new(sync.RWMutex),
	}
}

// The configuration for a single repository.
// Contains the information needed by the admission webhook to act on containers.
// If a container matches the regex specified as Name, then the policy specified will be applied.
type Repository struct {
	// The name of the repository, used to match images to policies.
	// Regexes are accepted (e.g. "registry/test/alpine.*", or "registry/.*")
	Name string `mapstructure:"name"`
	// Regex used to restrict to specific namespaces
	Namespace string `mapstructure:"namespace"`
	// Specifies the policy to be applied when the Name regex matches the container image.
	Trust    Trust `mapstructure:"trust"`
	Priority int   `mapstructure:"priority"`
}

// Created just to implement the priority sorting
type Repositories []Repository

// Needed methods to implement the Sort interface
func (r Repositories) Len() int {
	return len(r)
}

// Needed methods to implement the Sort interface
func (r Repositories) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Needed methods to implement the Sort interface
// Checks whether i should be before j?
func (r Repositories) Less(i, j int) bool {
	return r[i].Priority > r[j].Priority
}

// Specifies the behavior the webhook has to enforce on a matched container.
type Trust struct {
	// Whether Trust has to be enforced
	Enabled bool `mapstructure:"enabled,omitempty"`
	// The list of signers the matched image have to be tested against
	// If more than one signer is specified, images will have to be signed by all of them at the same time
	// in order to be accepted.
	Signers []*Signer `mapstructure:"signers,omitempty"`
	// The notary server that has to be used in order to verify an image
	TrustServer string `mapstructure:"trustServer,omitempty"`
}

// The info about a signer that has to be found
type Signer struct {
	// The Role name that has been used to sign the image.
	// Must be the exact role name used to sign the images on notary
	// e.g. "targets", "targets/releases" or "targets/username"
	Role string `mapstructure:"role"`
	// The public key of the keypair used to sign images in the repository
	PublicKey       string `mapstructure:"publicKey"`
	parsedPublicKey *data.PublicKey
	mutex           *sync.RWMutex
}

func (s *Signer) GetPEM(log *logrus.Entry) (*data.PublicKey, error) {
	s.mutex.RLock()
	if s.parsedPublicKey == nil {
		if err := s.parsePEM(log); err != nil {
			log.WithError(err).WithFields(logrus.Fields{"signer": s}).Error("Unable to parse signer public key")
			return nil, err
		}
	}
	defer s.mutex.RUnlock()
	return s.parsedPublicKey, nil
}

func (s *Signer) parsePEM(log *logrus.Entry) error {
	if s.parsedPublicKey == nil {
		log.WithField("signer", s.Role).Info("Initializing parsedPublicKey for signer")
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if s.parsedPublicKey == nil {
			log.WithField("signer", s).Debug("Parsing signature")
			pub, err := base64.StdEncoding.DecodeString(s.PublicKey)
			if err != nil {
				return err
			}
			keyFromConfig, err := utils.ParsePEMPublicKey(pub)
			if err != nil {
				return err
			}
			s.parsedPublicKey = &keyFromConfig
		}
		return nil

	}
	log.WithField("signer", s).WithField("publicKey", s.parsedPublicKey).Debug("Parsed signature")
	return nil
}

// given an image filter out the matching repositories
func (gc *GlobalConfig) GetMatchingRepositoriesPerImage(image, namespace string, log *logrus.Entry) ([]Repository, error) {
	c := gc.GetConfig()
	c.SortRepositories()
	repos := Repositories{}
	contextLogger := log.WithField("image", image)
	contextLogger.WithField("repositories", c.Repositories).Debug("searching for matching repos for image")
	atLeastOneNamespaceMatched := false
	for _, repo := range c.Repositories {
		matchedNamespace, err := regexp.Match(repo.Namespace, []byte(namespace))
		if err != nil {
			contextLogger.WithError(err).WithField("repo", repo).Error("namespace regex error")
		}
		if matchedNamespace {
			atLeastOneNamespaceMatched = true
			matched, err := regexp.Match(repo.Name, []byte(image))
			if err != nil {
				contextLogger.WithError(err).Error("Error matching repo regex to image")
			}
			if matched {
				contextLogger.Debug("Adding repo to returned repos because matched image")
				repos = append(repos, repo)
			}
		}
	}
	if !atLeastOneNamespaceMatched {
		contextLogger.Debug("no namespace matched")
		return nil, ErrNoNamespaceMatched{}
	} else if atLeastOneNamespaceMatched && len(repos) <= 0 {
		contextLogger.WithField("namespace", namespace).Debug("no repositories matched")
		return nil, ErrNoRepositoryMatched{}
	}
	contextLogger.WithField("repos", repos).Debug("Returning matched repositories")
	return repos, nil
}

func (c *Config) SortRepositories() {
	if !c.sortedRepositories {
		c.repositoriesMutex.Lock()
		defer c.repositoriesMutex.Unlock()
		if !c.sortedRepositories {
			sort.Sort(c.Repositories)
			c.sortedRepositories = true
		}
	}
}

func (c *Config) Validate(log *logrus.Entry) error {
	for _, repo := range c.Repositories {
		for _, signer := range repo.Trust.Signers {
			//TODO move mutex initialization to NewSigner function or something
			if signer.mutex == nil {
				log.WithField("signer", signer.Role).Info("Initializing mutex for signer")
				signer.mutex = new(sync.RWMutex)
			}
			err := signer.parsePEM(log)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"signer": signer, "repo": repo}).Error("Unable to parse signer public key")
				return err
			}
		}
	}
	return nil
}
