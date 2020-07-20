package reference

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
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
	//imageRegex  = regexp.MustCompile("^([a-zA-Z0-9-]+(.|:|/){0,1})+(@sha256:([a-fA-F0-9]+)){0,1}$")
)

// NewReference parses the image name and returns an error if the name is invalid.
func NewReference(name string, log *logrus.Entry) (*Reference, error) {
	//if !imageRegex.MatchString(name) {
	//return nil, errors.New("unsupported image")
	//}
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

	log.WithFields(logrus.Fields{"name": name, "reference": reference}).Debug("Obtained reference from name")

	return reference, nil
}

func (r *Reference) GetName() (name string) {
	name = r.Name
	if r.Tag != "" {
		name = fmt.Sprintf("%s:%s", name, r.Tag)
	}
	if r.Digest != "" {
		name = fmt.Sprintf("%s@sha256:%s", name, r.Digest)
	}
	return name
}
