package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sighupio/opa-notary-connector/internal/notary"
	"github.com/sighupio/opa-notary-connector/pkg/reference"
	"github.com/sirupsen/logrus"
)

const (
	// UUIDField is the key to be used in the context and logs to retrieve the call's uuid
	UUIDField = "uuid"
)

// Request contains the image that has to be checked against notary
type Request struct {
	Image string `json:"image,omitempty"`
}

// Response is returned when a Request is received
type Response struct {
	// Sha256 is the sha256 retrieved from notary.
	// Can be "", if from config the image doesn't have to be checked or if the image has not been signed by the required signers.
	// Check OK to know the reason.
	Sha256 string `json:"sha256,omitempty"`
	// Image is the complete image name (e.g. "<image_name>[:<tag>]@sha256:<sha256>").
	// Can be "", for the same reasons of Sha256
	Image string `json:"image,omitempty"`
	// OK is true if there has been no error, the image requested is fine to be deployed,
	// because signed by the required signers or because it did not require any
	OK bool `json:"ok"`
	// Err is the error message in case of OK set to false
	Err string `json:"error,omitempty"`
}

func GetImageHandlerBuilder(gc *conf.GlobalConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		log := logrus.WithField(UUIDField, c.GetString(UUIDField))
		image := c.Param("image")
		config := gc.GetConfig()

		response := Response{}

		sha256, image, err := CheckImage(image, config, gc.TrustRootDir, log)

		if err != nil {
			log.WithError(err).Errorf("there was an error while processing %+v", image)
			response.OK = false
			response.Err = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		response.OK = true
		response.Image = image
		response.Sha256 = sha256

		c.JSON(http.StatusOK, response)
	}

}

// CheckImageHandlerBuilder builds the main gin handler which will take a Request as input and return a Response
func CheckImageHandlerBuilder(gc *conf.GlobalConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		log := logrus.WithField(UUIDField, c.GetString(UUIDField))
		request := new(Request)
		if err := c.ShouldBindJSON(request); err != nil {
			log.WithError(err).Error("unable to bind body to Request object")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		config := gc.GetConfig()

		response := Response{}

		sha256, image, err := CheckImage(request.Image, config, gc.TrustRootDir, log)

		if err != nil {
			log.WithError(err).Errorf("there was an error while processing %+v", request)
			response.OK = false
			response.Err = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		response.OK = true
		response.Image = image
		response.Sha256 = sha256

		c.JSON(http.StatusOK, response)
	}

}

// CheckImage given an image and the config checks whether the image needs to be validated
// or not and if needed returns the sha and the complete image name (with sha)
func CheckImage(image string, config *conf.Config, trustRootDir string, log *logrus.Entry) (sha string, completeImage string, err error) {
	log = log.WithField("image", image)

	ref, _ := reference.NewReference(image, log)
	log.WithField("ref", ref).Debug("Got reference")

	repos, err := config.GetMatchingRepositoriesPerImage(ref, log)

	// if no repository matched, default deny and send
	if err != nil {
		log.WithError(err).Error("Got error when getting matching repositories")
		return "", "", err
	}

	log.WithField("repos", repos).Debug("Got matching repos for image")

	// repos are sorted by priority, therefore the first to be matched is the one with highest priority,
	// no other repos should be checked
	for _, repo := range repos {
		repo := repo
		// if one of the repos has no trust enabled and matches the image we should allow it
		if !repo.Trust.Enabled {
			return "", "", nil
		}
		no, err := notary.New(ref, &repo, trustRootDir, log)

		if err != nil {
			log.WithField("server", repo.Trust.TrustServer).WithError(err).Error("Not able to create cached repository for image")
			return "", "", err
		}

		// otherwise retrieve the signed sha from the repository and add the patch
		sha, err := no.GetSha()
		if err != nil {
			log.WithError(err).Error("Not able to get sha for image")
			return "", "", err
		}
		ref.Digest = sha
		return sha, ref.GetName(), nil

	}
	return "", "", nil
}
