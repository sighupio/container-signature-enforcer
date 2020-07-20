package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// Custom logging "middleware" to comply with the logging library used.
// Injects also the uuid that will be then injected in all
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// generating the call unique id
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			id, err := uuid.NewV4()
			if err != nil {
				logrus.WithError(err).Error("unable to generate uuid")
			}
			requestID = id.String()
		}
		c.Set(UUIDField, requestID)

		path := c.Request.URL.Path
		if query := c.Request.URL.RawQuery; query != "" {
			path = fmt.Sprintf("%s?%s", path, query)
		}
		requestLogger := logrus.WithFields(
			logrus.Fields{
				"state":      "received",
				"method":     c.Request.Method,
				"path":       path,
				"ip":         c.ClientIP(),
				"user-agent": c.Request.UserAgent(),
				UUIDField:    requestID,
			})
		requestLogger.Info("Request Received")

		c.Next()

		diff := time.Now().Sub(start)
		requestLogger.WithFields(
			logrus.Fields{
				"state":   "processed",
				"elapsed": diff,
				"result":  c.Writer.Status(),
			},
		).Info("Request Processed")
	}
}

func recoveryLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := logrus.WithField(UUIDField, c.GetString(UUIDField))
				log.WithFields(logrus.Fields{
					"error":      err,
					"stacktrace": strings.Split(string(debug.Stack()), "\n"),
				}).Error("Recovered panic")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
