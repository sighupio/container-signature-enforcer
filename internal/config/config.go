package config

import (
	"errors"
	"regexp"
	"sort"

	"github.com/sighupio/opa-notary-connector/pkg/reference"
	"github.com/sirupsen/logrus"
)

// Config is the main config structure containing informations regarding the image signature verification
type Config struct {
	Repositories `mapstructure:"repositories"`
	validated    bool
	//TODO: global signers that can be referenced by name from Trust, so that keys can be shared between repositories if needed
	//GlobalSigners []Signer
}

// GetMatchingRepositoriesPerImage given an image finds and returns only the matching repositories
func (c *Config) GetMatchingRepositoriesPerImage(image *reference.Reference, log *logrus.Entry) (repos Repositories, err error) {
	contextLogger := log.WithField("image", image)
	contextLogger.WithField("repositories", c.Repositories).Debug("searching for matching repos for image")
	for _, repo := range c.Repositories {
		matched, err := regexp.MatchString(repo.Name, image.Name)
		if err != nil {
			contextLogger.WithError(err).Error("Error matching repo regex to image")
		}
		if matched {
			contextLogger.Debug("Adding repo to returned repos because matched image")
			repos = append(repos, repo)
		}
	}
	if len(repos) == 0 {
		return nil, ErrNoRepositoryMatched{}
	}
	contextLogger.WithField("repos", repos).Debug("Returning matched repositories")
	return repos, nil
}

// Validate checks the config to be valid, mostly checks signers public keys are valid PEMs.
// It's expected to be run once and will return an error otherwise.
func (c *Config) Validate(log *logrus.Entry) error {
	if !c.validated {
		for _, repo := range c.Repositories {
			for _, signer := range repo.Trust.Signers {
				//TODO move mutex initialization to NewSigner function or something
				if err := signer.ParsePEM(log); err != nil {
					log.WithError(err).WithFields(logrus.Fields{"signer": signer, "repo": repo}).Error("Unable to parse signer public key")
					return err
				}
			}
		}
		sort.Sort(c.Repositories)
		c.validated = true
		return nil
	}
	return errors.New("Already validated config")
}
