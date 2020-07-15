package reference

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestReferenceOK(t *testing.T) {
	var tests = []struct {
		image, original, name, tag, digest, hostname, port string
	}{
		{image: "docker.io/test/alpine:test", original: "docker.io/test/alpine:test", name: "docker.io/test/alpine", tag: "test", digest: "", hostname: "docker.io", port: ""},
		{image: "registry.hub.docker.com/library/alpine:test", original: "registry.hub.docker.com/library/alpine:test", name: "registry.hub.docker.com/library/alpine", tag: "test", digest: "", hostname: "registry.hub.docker.com", port: ""},
		{image: "alpine", original: "alpine", name: "docker.io/library/alpine", tag: "latest", digest: "", hostname: "docker.io", port: ""},
		{image: "alpine@sha256:2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14", original: "alpine@sha256:2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14", name: "docker.io/library/alpine", tag: "latest", digest: "2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14", hostname: "docker.io", port: ""},
		{image: "registry.hub.docker.com:8080/library/alpine:test", original: "registry.hub.docker.com:8080/library/alpine:test", name: "registry.hub.docker.com:8080/library/alpine", tag: "test", digest: "", hostname: "registry.hub.docker.com", port: "8080"},
		{image: "localhost:30001/alpine:3.10", original: "localhost:30001/alpine:3.10", name: "localhost:30001/alpine", tag: "3.10", digest: "", hostname: "localhost", port: "30001"},
		{image: "localhost/alpine:3.10", original: "localhost/alpine:3.10", name: "localhost/alpine", tag: "3.10", digest: "", hostname: "localhost", port: ""},
	}
	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			ref, err := NewReference(tt.image, logrus.NewEntry(logrus.StandardLogger()))
			if err != nil {
				t.Errorf("Got error %s", err.Error())
				return
			}
			if ref == nil {
				t.Errorf("Got nil ref for %s", tt.original)
				return
			}
			if ref.Original != tt.original {
				t.Errorf("wanted %s, got %s as original", tt.original, ref.Original)
			}
			if ref.Name != tt.name {
				t.Errorf("wanted %s, got %s as name", tt.name, ref.Name)
			}
			if ref.Tag != tt.tag {
				t.Errorf("wanted %s, got %s as tag", tt.tag, ref.Tag)
			}
			if ref.Digest != tt.digest {
				t.Errorf("wanted %s, got %s as digest", tt.digest, ref.Digest)
			}
			if ref.Hostname != tt.hostname {
				t.Errorf("wanted %s, got %s as hostname", tt.hostname, ref.Hostname)
			}
			if ref.Port != tt.port {
				t.Errorf("wanted %s, got %s as port", tt.port, ref.Port)
			}
		})
	}
}

func TestReferenceName(t *testing.T) {
	var tests = []struct {
		image, expectedName string
	}{
		{image: "docker.io/test/alpine:test", expectedName: "docker.io/test/alpine:test"},
		{image: "registry.hub.docker.com/library/alpine:test", expectedName: "registry.hub.docker.com/library/alpine:test"},
		{image: "alpine", expectedName: "docker.io/library/alpine:latest"},
		{image: "alpine@sha256:2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14", expectedName: "docker.io/library/alpine:latest@sha256:2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14"},
		{image: "localhost:30001/alpine:3.10", expectedName: "localhost:30001/alpine:3.10"},
		{image: "localhost/alpine:3.10", expectedName: "localhost/alpine:3.10"},
	}
	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			ref, err := NewReference(tt.image, logrus.NewEntry(logrus.StandardLogger()))
			if err != nil {
				t.Errorf("Got error %s", err.Error())
				return
			}
			if ref == nil {
				t.Errorf("Got nil ref for %s", tt.image)
				return
			}
			if name := ref.GetName(); tt.expectedName != name {
				t.Errorf("Got %s expected %s as name", name, tt.expectedName)
			}
		})
	}
}

func TestMalformedImage(t *testing.T) {
	tests := []struct {
		image       string
		expectedRef Reference
	}{
		{image: "alpine:alksdja/asdasd:---", expectedRef: Reference{Original: "alpine:alksdja/asdasd:---", Name: "alpine:alksdja/asdasd", Tag: "---", Digest: "", Hostname: "", Port: ""}},
		{image: "alpine:alksdja/asdasd:/---", expectedRef: Reference{Original: "alpine:alksdja/asdasd:/---", Name: "alpine:alksdja/asdasd:/---", Tag: "latest", Digest: "", Hostname: "", Port: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			ref, err := NewReference(tt.image, logrus.NewEntry(logrus.StandardLogger()))
			if err != nil {
				t.Errorf("got error %v", err)
				return
			}

			if tt.expectedRef != *ref {
				t.Errorf("got %#v, expected %#v", ref, tt.expectedRef)
			}
		})
	}
}
