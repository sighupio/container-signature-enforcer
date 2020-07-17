package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/sighupio/opa-notary-connector/config"
)

func SetupServer(globalConfig *conf.GlobalConfig) *gin.Engine {
	r := gin.New()
	r.Use(ginLogger())
	r.Use(recoveryLogger())
	//TODO move to customRecovery to log with logrus on panic, will be available in next gin release
	//r.Use(gin.CustomRecovery())
	r.POST("/checkImage", CheckImageHandlerBuilder(globalConfig))
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "this is fine")
	})

	return r
}
