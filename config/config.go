package config

import (
	"regexp"
	"sort"

	"github.com/sighupio/opa-notary-connector/reference"
	"github.com/sirupsen/logrus"
)

// The main config structure containing informations regarding the image signature verification
type Config struct {
	Repositories `mapstructure:"repositories"`
	//TODO: global signers that can be referenced by name from Trust, so that keys can be shared between repositories if needed
	//GlobalSigners []Signer
}

// given an image filter out the matching repositories
func (c *Config) GetMatchingRepositoriesPerImage(image *reference.Reference, log *logrus.Entry) ([]Repository, error) {
	repos := Repositories{}
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
	if len(repos) <= 0 {
		return nil, ErrNoRepositoryMatched{}
	}
	contextLogger.WithField("repos", repos).Debug("Returning matched repositories")
	return repos, nil
}

func (c *Config) Validate(log *logrus.Entry) error {
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
	return nil
}
