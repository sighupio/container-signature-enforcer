package notary

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"testing"

	"github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sighupio/opa-notary-connector/pkg/reference"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/theupdateframework/notary"
	"github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/tuf/data"
)

func genPubKey(t *testing.T) string {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	assert.NoError(t, err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	var buf bytes.Buffer
	err = pem.Encode(&buf, pemkey)
	assert.NoError(t, err)

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func TestRepository(t *testing.T) {
	t.Parallel()
	log := logrus.NewEntry(&logrus.Logger{})
	conf := config.Config{}
	signer := &config.Signer{
		Role:      "targets/sighup",
		PublicKey: genPubKey(t),
	}
	conf.Repositories = config.Repositories{
		config.Repository{
			//Priority is not checked here
			Priority: 11,
			Name:     "docker.io.*",
			Trust: &config.Trust{
				Enabled:     true,
				TrustServer: "sighup.notary.com",
				Signers: []*config.Signer{
					signer,
				},
			},
		},
	}
	// Validate also initialize mutexes in repository and signers => shit
	err := conf.Validate(log)
	assert.NoError(t, err)
	pubKey, err := signer.GetPEM(log)
	assert.NoError(t, err)
	var fakeMetadataGetter AllTargetMetadataByNameGetter = fake{
		map[string][]client.TargetSignedStruct{
			"not-latest": {
				client.TargetSignedStruct{
					Role: data.DelegationRole{
						BaseRole: data.BaseRole{
							Name: "targets/sighup",
							Keys: map[string]data.PublicKey{
								(*pubKey).ID(): *pubKey,
							},
						},
					},
					Target: client.Target{
						Hashes: data.Hashes{
							notary.SHA256: []byte("not-sighup"),
						},
					},
				},
				client.TargetSignedStruct{
					Role: data.DelegationRole{
						BaseRole: data.BaseRole{
							Name: "targets/releases",
							Keys: map[string]data.PublicKey{
								(*pubKey).ID(): *pubKey,
							},
						},
					},
					Target: client.Target{
						Hashes: data.Hashes{
							notary.SHA256: []byte("not-sighup"),
						},
					},
				},
			},
			"latest": {
				client.TargetSignedStruct{
					Role: data.DelegationRole{
						BaseRole: data.BaseRole{
							Name: "targets/sighup",
							Keys: map[string]data.PublicKey{
								(*pubKey).ID(): *pubKey,
							},
						},
					},
					Target: client.Target{
						Hashes: data.Hashes{
							notary.SHA256: []byte("sighup"),
						},
					},
				},
			},
		},
	}

	var tests = []struct {
		image              string
		repo               *config.Repository
		fakeMetadataGetter *AllTargetMetadataByNameGetter
		expectedSha        string
	}{
		{image: "docker.io/library/alpine", fakeMetadataGetter: &fakeMetadataGetter, repo: &conf.Repositories[0], expectedSha: "sighup"},
		{image: "docker.io:8080/library/alpine", fakeMetadataGetter: &fakeMetadataGetter, repo: &conf.Repositories[0], expectedSha: "sighup"},
		{image: "alpine:not-latest", fakeMetadataGetter: &fakeMetadataGetter, repo: &conf.Repositories[0], expectedSha: "not-sighup"},
		{image: "alpine:latest", fakeMetadataGetter: &fakeMetadataGetter, repo: &conf.Repositories[0], expectedSha: "sighup"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.image, func(t *testing.T) {
			t.Parallel()
			ref, _ := reference.NewReference(tt.image, logrus.NewEntry(logrus.StandardLogger()))
			repo, err := newWithGetter(ref, tt.repo, tt.fakeMetadataGetter, "", log)
			assert.NoError(t, err)

			encodedExpectedSha := hex.EncodeToString([]byte(tt.expectedSha))
			sha, err := repo.GetSha()
			assert.NoError(t, err)
			assert.Equal(t, encodedExpectedSha, sha)

		})
	}
}
