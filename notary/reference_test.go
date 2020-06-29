package notary

import (
	"testing"
)

func TestReferenceOK(t *testing.T) {
	var tests = []struct {
		image, original, name, tag, digest, hostname, port string
	}{
		{image: "docker.io/test/alpine:test", original: "docker.io/test/alpine:test", name: "docker.io/test/alpine", tag: "test", digest: "", hostname: "docker.io", port: ""},
		{image: "registry.hub.docker.com/library/alpine:test", original: "registry.hub.docker.com/library/alpine:test", name: "registry.hub.docker.com/library/alpine", tag: "test", digest: "", hostname: "registry.hub.docker.com", port: ""},
		{image: "alpine", original: "alpine", name: "docker.io/library/alpine", tag: "latest", digest: "", hostname: "docker.io", port: ""},
		{image: "alpine@sha256:2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14", original: "alpine@sha256:2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14", name: "docker.io/library/alpine", tag: "latest", digest: "2bb501e6173d9d006e56de5bce2720eb06396803300fe1687b58a7ff32bf4c14", hostname: "docker.io", port: ""},
	}
	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			ref, err := NewReference(tt.image)
			if err != nil {
				t.Errorf("Got error %s", err.Error())
			}
			if ref == nil {
				t.Errorf("Got nil ref for %s", tt.original)
				return
			}
			if ref.original != tt.original {
				t.Errorf("wanted %s, got %s as original", tt.original, ref.original)
			}
			if ref.name != tt.name {
				t.Errorf("wanted %s, got %s as name", tt.name, ref.name)
			}
			if ref.tag != tt.tag {
				t.Errorf("wanted %s, got %s as tag", tt.tag, ref.tag)
			}
			if ref.digest != tt.digest {
				t.Errorf("wanted %s, got %s as digest", tt.digest, ref.digest)
			}
			if ref.hostname != tt.hostname {
				t.Errorf("wanted %s, got %s as hostname", tt.hostname, ref.hostname)
			}
			if ref.port != tt.port {
				t.Errorf("wanted %s, got %s as port", tt.port, ref.port)
			}
		})
	}
}
