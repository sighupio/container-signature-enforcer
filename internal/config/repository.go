package config

// The configuration for a single repository.
// Contains the information needed by the admission webhook to act on containers.
// If a container matches the regex specified as Name, then the policy specified will be applied.
type Repository struct {
	// The name of the repository, used to match images to policies.
	// Regexes are accepted (e.g. "registry/test/alpine.*", or "registry/.*")
	Name string `mapstructure:"name"`
	// Specifies the policy to be applied when the Name regex matches the container image.
	Trust    *Trust `mapstructure:"trust"`
	Priority int    `mapstructure:"priority"`
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

type Credentials struct {
	//Secret   string `mapstructure:"secret,omitempty"`
	User     string `mapstructure:"user,omitempty"`
	Password string `mapstructure:"pass,omitempty"`
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
	TrustServer string       `mapstructure:"trustServer,omitempty"`
	Credentials *Credentials `mapstructure:"auth,omitempty"`
}

func (c *Credentials) GetCreds() (string, string, error) {
	// could retrieve creds from a secret at some point
	return c.User, c.Password, nil
}
