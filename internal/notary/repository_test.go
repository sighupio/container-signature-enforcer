package notary

import (
	"encoding/hex"
	"testing"

	"github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sighupio/opa-notary-connector/pkg/reference"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/theupdateframework/notary"
	"github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/tuf/data"
)

func TestRepository(t *testing.T) {
	t.Parallel()
	log := logrus.NewEntry(&logrus.Logger{})
	conf := config.Config{}
	signer := &config.Signer{
		Role:      "targets/sighup",
		PublicKey: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURhekNDQWxPZ0F3SUJBZ0lVYXQzbjRucm5IUGxrdHdJR2F5ZGxqaVFpUFRNd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TURFeU1ERXhORE0xTVRCYUZ3MHlNVEV5Ck1ERXhORE0xTVRCYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdnZ0VpTUEwR0NTcUdTSWIzRFFFQgpBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQzE5bDdTeEpSSUhVbkV5TTFNeU1TcUFYS1Y4YWVuRDR5aExuNlROcjVFCmxCaUhtZytIU3F2KytSQUV1NmhoRU5TU2t2cFd3aGJ3U0ZTU2ZTS0tDZm5WNHdBRElvekxRZlFPZUMyc2NhM0UKSGNwNkFYQ1FtSGRrYWx2NHl4SEhEazEweDI0TDVTTmF3ZFYvUm5DQUJ5VUxmeUs5QkJSYzhDQVlEQmpaS1R4bwpJM2J6UlBUVFNSd1VIWHJFMDhpYmkxS2lsTFprT2hKWGhYd0ltbjhmbkJtTHdTd0EwQWRIazF2cWFQR2J4bnVKCmU1U0JiVEljdjBmOEFCaE5rQ3ZRWVlyTEtrM2x5NjhWMktDK0ZTRmNoangrUFdRTUVTOWZ6NUZUU3F4YUY2ZHkKSzR4MFRMbWpNcDBabkowWU1qYlFaK1I3SXFtRy83TzhvVlJZQXlrdEZyZ2hBZ01CQUFHalV6QlJNQjBHQTFVZApEZ1FXQkJUKzd2dFNOREZpWW9YRDc2Y0dVMW9TdXk0M3F6QWZCZ05WSFNNRUdEQVdnQlQrN3Z0U05ERmlZb1hECjc2Y0dVMW9TdXk0M3F6QVBCZ05WSFJNQkFmOEVCVEFEQVFIL01BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQ2EKM1JzUUpIazNTUWlTMCtCdEFuSGNGZ09scDU0OGcrMWZlcXFuWmdYZWMvUHhsSU1tUXVvSHExNlNTTUMxcFJTagphM1liMy9uZUkvR3d6ZzlIc1pRYm0vdlA0YjlNaG55aXN5ZGxacHhKS0ZIdWExN1hjcWV0VUJBUEpFUC9LVUVjClNNUDA4dXlETURKNVcxcVdZOFNnczUxeEM2Z1cycTZJd3d5Q3gvM25EVy9PSXhKc0NEYVJQNHp1WmRmMWw2ZGIKYmhUMWY0YUlyRUlvSVlPeFhOV05tQ1FMUWFiMjVncWw2aVJ5TVFieGx6a0kxTm43NWJsbGFyZjJ4QWpLSmNHRQpVQWFKNkZwWjhENGRkQnhCZmF3L29oNld1UXVscmlFWEozc05DMndLWWtHdnlTZi9VcmZQb2RKTEFvKy9BY3BSCnZZM0d1OFl6dUMwZE9WUkFtL3RPCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
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
