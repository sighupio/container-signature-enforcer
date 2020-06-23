package cmd

import (
	"fmt"
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
		trxID, err := uuid.NewV4()
		if err != nil {
			logrus.WithError(err).Error("unable to generate uuid")
		}
		c.Set(uuidField, trxID.String())

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
				uuidField:    trxID,
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

type recoveryLogger struct {
}

func (l recoveryLogger) Write(p []byte) (n int, err error) {
	logrus.WithField("panic", true).Error(strings.Split("\n", string(p)))
	return len(p), nil
}
