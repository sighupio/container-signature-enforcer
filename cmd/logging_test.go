package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoveryLogging(t *testing.T) {
	r := setupServer()
	r.GET("/panic", func(c *gin.Context) {
		panic("testing recovery logger")
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/panic", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("got wrong status code")
	}

}
