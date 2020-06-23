package admission

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sighupio/opa-notary-connector/config"
	conf "github.com/sighupio/opa-notary-connector/config"
	"github.com/sighupio/opa-notary-connector/notary"
	logrus "github.com/sirupsen/logrus"
	vAdmission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//const (
//AnnotationsPath         = "/metadata/annotations"
//AcceptedAnnotationLabel = "opa-notary-connector-accepted"
//AcceptedAnnotationPath  = AnnotationsPath + "/" + AcceptedAnnotationLabel
//)

type JSONPatch struct {
	From  string      `json:"from,omitempty"`
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// Given an admission request having a pod as subject + the necessary configurations, this function returns a handler with the following logic:
// 1. umarshal the admissionReview request, return StatusBadRequest if error
// 2. extract pod object, return StatusBadRequest if error
// 3. prepare the AdmissionResponse
// 4. retrieve all containers/initcontainers
// 5. foreach container check against matching policies:
//		5a. if trust is disabled for one of the matched policies, continue with other containers else check against provided notary server, retrieve the latest signed sha and store the patch to modify the pod object.
//		5b. first deny => deny the request, otherwise accumulate all patches
func HandleAdmissionRequest(config *conf.GlobalConfig) func(c *gin.Context) {
	return func(c *gin.Context) {

		log := logrus.WithField("uuid", c.GetString("uuid"))

		ar, err := getAdmissionReview(c, log)

		if err != nil {
			log.WithError(err).Errorf("unable to extract admissionReview")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		status, ar, err := reviewAdmission(ar, config, log)

		if err != nil {
			c.AbortWithError(status, err)
		}
		c.JSON(status, ar)
	}
}

func getContainerSpecPathMap(podSpec *corev1.PodSpec, baseSpecPath string, log *logrus.Entry) map[string][]string {
	containerSpecPath := map[string][]string{}
	for i, container := range podSpec.Containers {
		path := fmt.Sprintf("%s/containers/%d/image", baseSpecPath, i)
		containerSpecPath[container.Image] = append(containerSpecPath[container.Image], path)
		log.WithFields(logrus.Fields{"image": container.Image, "index": i, "path": path}).Debug("Retrieved container")
	}
	for i, initcontainer := range podSpec.InitContainers {
		path := fmt.Sprintf("%s/initContainers/%d/image", baseSpecPath, i)
		containerSpecPath[initcontainer.Image] = append(containerSpecPath[initcontainer.Image], path)
		log.WithFields(logrus.Fields{"image": initcontainer.Image, "index": i, "path": path}).Debug("Retrieved initContainer")
	}
	return containerSpecPath
}

func reviewAdmission(ar *vAdmission.AdmissionReview, config *config.GlobalConfig, log *logrus.Entry) (int, *vAdmission.AdmissionReview, error) {
	namespace := ar.Request.Namespace

	baseSpecPath, podSpec, err := getPodSpec(ar, log)
	if err != nil {
		log.WithError(err).Error("Error getting podSpec from request")
		return http.StatusBadRequest, nil, err
	}
	log.WithFields(logrus.Fields{"baseSpecPath": baseSpecPath, "podSpec": podSpec}).Debug("Retrieved podSpec from request")

	// 3. initialize the response, default deny
	admissionResponse := &vAdmission.AdmissionResponse{Allowed: false}

	// 4. extract all images from containers and initcontainers
	// store also their path in the spec for later patching
	ContainerSpecPath := getContainerSpecPathMap(podSpec, baseSpecPath, log)

	// 5. check all images against policies defined in the config
	patches := []JSONPatch{}
	status := -1
	for image, specPaths := range ContainerSpecPath {
		status, err := admissionLogic(admissionResponse, ar, namespace, image, specPaths, patches, log, config)
		if err != nil {
			admissionResponse = denyAdmission(admissionResponse, fmt.Sprintf("The image %s is not allowed, given that is not matching any configured repository", image), log)
			status, ar = prepareResponse(admissionResponse, patches, ar, log)
			return status, ar, err
		}
	}
	// FIN QUA
	status, ar = prepareResponse(admissionResponse, patches, ar, log)
	return status, ar, err
}

//TODO:
// namespace input param will be moved to rego
// wrap it in handler for body { "image": "...", "namespace": "..." }
// inject Notary dependency to make it testable
func referee(namespace, image string, log *logrus.Entry, config *conf.GlobalConfig) (ok bool, imageWithSha string, err error) {
	repos, err := config.GetMatchingRepositoriesPerImage(strings.Split(image, ":")[0], namespace, log)
	log.WithFields(logrus.Fields{"image": image, "repos": repos}).Debug("Got matching repos for image")

	// if no repository matched, default deny and send
	if err != nil {
		return false, "", err
	}

	// repos are sorted by priority, therefore the first to be matched is the one with highest priority,
	// no other repos should be checked
	for _, repo := range repos {
		// if one of the repos has no trust enabled and matches the image we should allow it
		if !repo.Trust.Enabled {
			return true, "", nil
		} else {
			ref, err := notary.NewReference(image)

			if err != nil {
				log.WithFields(logrus.Fields{
					"image":  image,
					"server": repo.Trust.TrustServer,
				}).WithError(err).Error("Image was not parsable")
				return false, "", err
			}

			client, err := notary.NewFileCachedRepository(config, &repo, ref, log)

			if err != nil {
				log.WithFields(logrus.Fields{
					"image":  image,
					"server": repo.Trust.TrustServer,
				}).WithError(err).Error("Not able to create cached repository for image")
				return false, "", err
			}

			// otherwise retrieve the signed sha from the repository and add the patch
			imageWithSha, err := notary.CheckImage(ref, config.TrustRootDir, &repo, client, log)
			return true, imageWithSha, err

		}
	}
	return false, "", err
}

func refereeLoop(namespace, image string, specPaths []string, patches []JSONPatch, log *logrus.Entry, config *conf.GlobalConfig) (bool, error) {
	ok, imageWithSha, err := referee(namespace, image, log, config)
	if err != nil {
		return false, err
	} else {
		for _, specPath := range specPaths {
			patch := JSONPatch{
				Op:    "replace",
				Path:  specPath,
				Value: imageWithSha,
			}
			log.WithField("patch", patch).Debug("Adding single patch")
			patches = append(patches, patch)
		}
		log.WithField("patches", patches).Debug("Patches to be returned")
	}
	return ok, err

}

// TODO: remove admissionResponse, ar, specPaths
func admissionLogic(admissionResponse *vAdmission.AdmissionResponse, ar *vAdmission.AdmissionReview, namespace, image string, specPaths []string, patches []JSONPatch, log *logrus.Entry, config *config.GlobalConfig) (int, error) {
	// get repositories matching the image

	status := -1
	ok, err := refereeLoop(namespace, image, specPaths, patches, log, config)
	// if no repository matched, default deny and send
	if err != nil {
		contextLogger := log.WithFields(logrus.Fields{"namespace": namespace, "image": image})
		switch err := err.(type) {
		case conf.ErrNoNamespaceMatched:
			contextLogger.Debug("No namespace matched, will allow")
			admissionResponse = allowAdmission(admissionResponse, log)
		case conf.ErrNoRepositoryMatched:
			contextLogger.Debug("No repos matching, will deny")
			admissionResponse = denyAdmission(admissionResponse, fmt.Sprintf("The image %s is not allowed, given that is not matching any configured repository", image), log)
		default:
			contextLogger.WithError(err).Error("Unkown error returned while trying to get matching repositories")
			admissionResponse = denyAdmission(admissionResponse, fmt.Sprintf("Unexpected error: %v", err), log)
		}

		status, _ = prepareResponse(admissionResponse, patches, ar, log)
	}
	if ok {
		admissionResponse = allowAdmission(admissionResponse, log)
	} else {
		admissionResponse = denyAdmission(admissionResponse, "generic error", log)
	}
	status, _ = prepareResponse(admissionResponse, patches, ar, log)
	return status, err

}

func prepareResponse(admissionResponse *vAdmission.AdmissionResponse, patches []JSONPatch, reviewRequest *vAdmission.AdmissionReview, log *logrus.Entry) (int, *vAdmission.AdmissionReview) {
	//patches = append(patches, JSONPatch{
	//Op:    "add",
	//Path:  AcceptedAnnotationPath,
	//Value: map[string]"true",
	//})

	if len(patches) > 0 {
		log.WithField("patches", patches).Debug("Sending patches")
		jsonPatch, err := json.Marshal(patches)
		if err != nil {
			log.WithError(err).Error("error marshalling patches to json")
			denyAdmission(admissionResponse, err.Error(), log)
		} else {
			pt := vAdmission.PatchTypeJSONPatch
			admissionResponse.PatchType = &pt
			admissionResponse.Patch = jsonPatch
		}
	}

	reviewRequest.Response = admissionResponse

	log.WithField("response", reviewRequest).Debug("Sending response")
	log.WithFields(logrus.Fields{"patches": patches}).Info("Sending response")

	return http.StatusOK, reviewRequest
}

func allowAdmission(admissionResponse *vAdmission.AdmissionResponse, log *logrus.Entry) *vAdmission.AdmissionResponse {
	admissionResponse.Allowed = true
	admissionResponse.Result = nil
	log.WithField("response", admissionResponse).Debug("Allowing")
	return admissionResponse
}

func denyAdmission(admissionResponse *vAdmission.AdmissionResponse, message string, log *logrus.Entry) *vAdmission.AdmissionResponse {
	admissionResponse.Allowed = false
	admissionResponse.Result = &metav1.Status{
		Status:  message,
		Message: message,
		Code:    http.StatusBadRequest,
		Reason:  metav1.StatusReasonBadRequest,
	}
	log.WithField("response", admissionResponse).Debug("Denying")
	return admissionResponse
}
