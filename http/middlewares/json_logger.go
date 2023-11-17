package middlewares

import (
	"time"

	util "github.com/OdairPianta/gryphon_api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// JSONLogMiddleware logs a gin HTTP request in JSON format, with some additional custom key/values
func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := util.GetDurationInMillseconds(start)

		entry := log.WithFields(log.Fields{
			"client_ip": util.GetClientIP(c),
			"duration":  duration,
			"method":    c.Request.Method,
			"path":      c.Request.RequestURI,
			"status":    c.Writer.Status(),
			"body":      c.Request.Body,
			// "api_version": util.ApiVersion,
		})

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("")
		}
	}
}
