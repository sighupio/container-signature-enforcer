package reference

import (
	"regexp"
	"strings"
)

// Reference .
type Reference struct {
	Original string
	Name     string
	Tag      string
	Digest   string
	Hostname string
	Port     string
}

var (
	digestRegex = regexp.MustCompile("@sha256:(?P<sha256>[a-fA-F0-9]+)$")
	tagRegex    = regexp.MustCompile(":(?P<tag>[^/]+)$")
	hostRegex   = regexp.MustCompile("^(?P<host>[^/^:]*)(/|(:(?P<port>[0-9]+)))")
)

// NewReference parses the image name and returns an error if the name is invalid.
func NewReference(name string) (*Reference, error) {
	reference := &Reference{}
	reference.Original = name

	if !strings.Contains(name, "/") {
		name = "docker.io/library/" + name
	}

	if digestRegex.MatchString(name) {
		res := digestRegex.FindStringSubmatch(name)
		reference.Digest = res[1] // digest capture group index
		name = strings.TrimSuffix(name, res[0])
	}
	if tagRegex.MatchString(name) {
		res := tagRegex.FindStringSubmatch(name)
		reference.Tag = res[1] // tag capture group index
		name = strings.TrimSuffix(name, res[0])
	} else {
		reference.Tag = "latest"
	}

	// everything else is the name
	reference.Name = name

	if hostRegex.MatchString(name) {
		res := hostRegex.FindStringSubmatch(name)
		reference.Hostname = res[1] // host capture group index
		reference.Port = res[4]     // port capture group index, could be empty string if not matched
	}

	return reference, nil
}