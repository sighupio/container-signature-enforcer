package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/sighupio/opa-notary-connector/config"
	"github.com/sighupio/opa-notary-connector/notary"
	"github.com/sighupio/opa-notary-connector/reference"
	"github.com/sirupsen/logrus"
)

const (
	UUIDField = "uuid"
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

func CheckImageHandlerBuilder(gc *conf.GlobalConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		config := gc.GetConfig()

		response := Response{}
		log := logrus.WithField(UUIDField, c.GetString(UUIDField))

		request := new(Request)
		if err := c.ShouldBindJSON(request); err != nil {
			log.WithError(err).Error("unable to bind body to Request object")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		sha256, err := CheckImage(request.Image, config, gc.TrustRootDir, log)

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

func CheckImage(image string, config *conf.Config, trustRootDir string, log *logrus.Entry) (sha string, err error) {
	log = log.WithField("image", image)

	ref, _ := reference.NewReference(image, log)
	log.WithField("ref", ref).Debug("Got reference")

	repos, err := config.GetMatchingRepositoriesPerImage(ref, log)

	// if no repository matched, default deny and send
	if err != nil {
		log.WithError(err).Error("Got error when getting matching")
		return "", err
	}

	log.WithField("repos", repos).Debug("Got matching repos for image")

	// repos are sorted by priority, therefore the first to be matched is the one with highest priority,
	// no other repos should be checked
	for _, repo := range repos {
		// if one of the repos has no trust enabled and matches the image we should allow it
		if !repo.Trust.Enabled {
			return "", nil
		} else {

			no, err := notary.New(ref, &repo, trustRootDir, log)

			if err != nil {
				log.WithField("server", repo.Trust.TrustServer).WithError(err).Error("Not able to create cached repository for image")
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
