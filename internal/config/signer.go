package config

import (
	"encoding/base64"

	"github.com/sirupsen/logrus"
	"github.com/theupdateframework/notary/tuf/data"
	"github.com/theupdateframework/notary/tuf/utils"
)

// The info about a signer that has to be found
type Signer struct {
	// The Role name that has been used to sign the image.
	// Must be the exact role name used to sign the images on notary
	// e.g. "targets", "targets/releases" or "targets/username"
	Role string `mapstructure:"role"`
	// The public key of the keypair used to sign images in the repository
	PublicKey       string `mapstructure:"publicKey"`
	parsedPublicKey *data.PublicKey
}

func (s *Signer) GetPEM(log *logrus.Entry) (*data.PublicKey, error) {
	if s.parsedPublicKey == nil {
		if err := s.ParsePEM(log); err != nil {
			log.WithError(err).WithFields(logrus.Fields{"signer": s}).Error("Unable to parse signer public key")
			return nil, err
		}
	}
	return s.parsedPublicKey, nil
}

func (s *Signer) ParsePEM(log *logrus.Entry) error {
	if s.parsedPublicKey == nil {
		log.WithField("signer", s.Role).Info("Initializing parsedPublicKey for signer")
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
