package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Loosely based on:
// https://github.com/sbecker/gin-api-demo/blob/07f9a9242f743fc51ae4a046ee58e12627bad571/middleware/json_logger.go#L1
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		c.Header("X-Start", fmt.Sprintf("%.7f", float64(start.UnixNano())/float64(1000000000)))
		// Process Request
		c.Next()

		// Stop timer
		end := time.Now()
		duration := GetDurationInMillseconds(start, end)

		entry := log.WithFields(log.Fields{
			"duration": duration,
			"start":    start,
			"end":      end,
			"path":     c.Request.RequestURI,
			"status":   c.Writer.Status(),
		})

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("")
		}
	}
}

// Based on https://github.com/sbecker/gin-api-demo/blob/07f9a9242f743fc51ae4a046ee58e12627bad571/util/log.go#L86
// GetDurationInMillseconds takes a start time and returns a duration in milliseconds
func GetDurationInMillseconds(start time.Time, end time.Time) float64 {
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}
