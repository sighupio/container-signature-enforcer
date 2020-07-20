package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/sighupio/opa-notary-connector/internal/config"
)

// SetupServer configures the gin http server with:
// - the logrus request logger
// - the logrus recovery logger
// - the handlers required
func SetupServer(globalConfig *conf.GlobalConfig) *gin.Engine {
	r := gin.New()
	r.Use(ginLogger())
	r.Use(recoveryLogger())
	r.POST("/checkImage", CheckImageHandlerBuilder(globalConfig))
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "this is fine")
	})

	return r
}
