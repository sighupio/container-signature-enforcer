package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sighupio/opa-notary-connector/config"
	"github.com/stretchr/testify/assert"
)

func TestHealthz(t *testing.T) {
	t.Parallel()
	r := SetupServer(config.NewGlobalConfig())
	ts := httptest.NewServer(r)
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/healthz", ts.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
