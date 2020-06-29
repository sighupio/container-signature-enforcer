package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/sighupio/opa-notary-connector/config"
	"github.com/sighupio/opa-notary-connector/core"
	"github.com/sirupsen/logrus"
)

func ImageShaBuilder(config *conf.GlobalConfig) func(c *gin.Context) {
	return func(c *gin.Context) {

		log := logrus.WithField("uuid", c.GetString("uuid"))

		request := new(Request)
		if err := c.ShouldBindJSON(request); err != nil {
			log.WithError(err).Error("unable to bind body to Request object")
			c.AbortWithError(http.StatusBadRequest, err)
		}

		//TODO remove namespace
		sha256, err := core.Referee(request.Namespace, request.Image, log, config)

		if err != nil {
			log.WithError(err).Errorf("there was an error while processing %+v", request)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		response := Response{
			Request: *request,
			Sha256:  sha256,
		}

		c.JSON(http.StatusOK, response)
	}

}

type Request struct {
	Namespace string `json:"namespace,omitempty"`
	Image     string `json:"image,omitempty"`
}

type Response struct {
	Request
	Sha256 string `json:"sha256,omitempty"`
}
