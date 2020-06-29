package notary

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/distribution/reference"
)

// Reference .
type Reference struct {
	original string
	name     string
	tag      string
	digest   string
	hostname string
	port     string
}

// NewReference parses the image name and returns an error if the name is invalid.
func NewReference(name string) (*Reference, error) {
	var digest string
	original := name
	// Remove the digest so `ParseNamed` doesn't fail, it can't handle short digests.
	if strings.Contains(name, "@sha256:") {
		fields := strings.Split(name, "@sha256:")
		name = fields[0]
		digest = fields[1]
	}

	if !strings.Contains(name, "/") {
		name = fmt.Sprintf("docker.io/library/%s", name)
	}

	// Get image name
	ref, err := reference.ParseNamed(name)
	if err != nil {
		return nil, err
	}

	// Get the hostname
	hostname, _ := reference.SplitHostname(ref)
	if hostname == "" {
		// If no domain found, treat it as docker.io
		hostname = "docker.io"
	}
	if !strings.Contains(hostname, ".") {
		// Fix SplitHostname wrongly splitting repositories like molepigeon/wibble
		hostname = "docker.io"
	}
	// Make sure it can be used to build a valid URL
	u, err := url.Parse("http://" + hostname)
	if err != nil {
		return nil, err
	}

	// if the image does not have a tag, use `latest` so we can parse it again.
	image := strings.Replace(name, hostname, "", 1)
	if !strings.Contains(image, ":") {
		name += ":latest"
	}

	// Parse the name again including the tag so we can have a reference.taggedReference object
	// we ommit the error here since we already parsed the original string above.
	ref, _ = reference.ParseNamed(name)

	return &Reference{
		original: original,
		name:     ref.Name(),
		tag:      ref.(reference.Tagged).Tag(),
		digest:   digest,
		hostname: u.Hostname(),
		port:     u.Port(),
	}, nil
}
