package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/sighupio/opa-notary-connector/internal/config"
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

func TestFuzzingInput(t *testing.T) {
	t.Parallel()
	for i := 1; i <= 20; i++ {
		i := i
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			r := SetupServer(config.NewGlobalConfig())
			ts := httptest.NewServer(r)
			defer ts.Close()
			req := &Request{}
			f := fuzz.New()
			f.Fuzz(req)
			jsonReq, err := json.Marshal(req)
			assert.NoError(t, err)
			resp, err := http.Post(fmt.Sprintf("%s/checkImage", ts.URL), "application/json", bytes.NewBuffer(jsonReq))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			response := &Response{}
			defer resp.Body.Close()
			err = json.Unmarshal(body, response)
			assert.NoError(t, err)
			assert.NotEqual(t, response.Err, "")
			assert.Equal(t, response.Image, "")
			assert.Equal(t, response.Sha256, "")
		})
	}
}
