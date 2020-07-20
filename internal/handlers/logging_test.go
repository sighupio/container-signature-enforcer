package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryLogging(t *testing.T) {
	t.Parallel()
	logrus.SetFormatter(new(logrus.JSONFormatter))
	r := SetupServer(config.NewGlobalConfig())
	r.GET("/panic", func(c *gin.Context) {
		var nothing *struct{ test string }
		_ = fmt.Sprint(nothing.test)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/panic", ts.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestLogging(t *testing.T) {
	t.Parallel()
	logrus.SetFormatter(new(logrus.JSONFormatter))
	r := SetupServer(config.NewGlobalConfig())
	r.GET("/echo", func(c *gin.Context) {
		id := c.GetString(UUIDField)
		c.String(http.StatusOK, id)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/echo", ts.URL), nil)
	hardcodedID := "TEST-ID"
	req.Header.Add("X-Request-ID", hardcodedID)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, hardcodedID, string(body))

}
