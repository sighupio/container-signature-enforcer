package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	conf "github.com/sighupio/opa-notary-connector/config"
	"github.com/sighupio/opa-notary-connector/notary"
	"github.com/sirupsen/logrus"
)

type Request struct {
	Image string `json:"image,omitempty"`
}

type Response struct {
	Request
	Sha256 string `json:"sha256,omitempty"`
	OK     bool   `json:"ok"`
	Err    string `json:"error,omitempty"`
}

func CheckImageHandlerBuilder(config *conf.GlobalConfig) func(c *gin.Context) {
	return func(c *gin.Context) {

		response := Response{}
		log := logrus.WithField("uuid", c.GetString("uuid"))

		request := new(Request)
		if err := c.ShouldBindJSON(request); err != nil {
			log.WithError(err).Error("unable to bind body to Request object")
			c.AbortWithError(http.StatusBadRequest, err)
		}

		sha256, err := CheckImage(request.Image, log, config)

		if err != nil {
			log.WithError(err).Errorf("there was an error while processing %+v", request)
			response.OK = false
			response.Err = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		response.OK = true
		response.Request = *request
		response.Sha256 = sha256

		c.JSON(http.StatusOK, response)
	}

}

func CheckImage(image string, log *logrus.Entry, config *conf.GlobalConfig) (sha string, err error) {
	repos, err := config.GetMatchingRepositoriesPerImage(strings.Split(image, ":")[0], log)
	log.WithFields(logrus.Fields{"image": image, "repos": repos}).Debug("Got matching repos for image")

	// if no repository matched, default deny and send
	if err != nil {
		return "", err
	}

	// repos are sorted by priority, therefore the first to be matched is the one with highest priority,
	// no other repos should be checked
	for _, repo := range repos {
		// if one of the repos has no trust enabled and matches the image we should allow it
		if !repo.Trust.Enabled {
			return "", nil
		} else {

			no, err := notary.New(image, &repo, log)

			if err != nil {
				log.WithFields(logrus.Fields{
					"image":  image,
					"server": repo.Trust.TrustServer,
				}).WithError(err).Error("Not able to create cached repository for image")
				return "", err
			}

			// otherwise retrieve the signed sha from the repository and add the patch
			sha, err := no.GetSha()
			if err != nil {
				log.WithError(err).Error("Not able to get sha for image")
				return "", err
			}
			return sha, nil

		}
	}
	return "", err
}
