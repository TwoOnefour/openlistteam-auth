// Package utils utils/log.go
package utils

import (
	"github.com/gin-gonic/gin"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetLevel(logrus.InfoLevel)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		entry := logrus.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"duration":   duration.Seconds(),
		})
		c.Set("logger", entry)
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		}

	}
}

func GetLogger(c *gin.Context) *logrus.Entry {
	if le, ok := c.Get("logger"); ok {
		if entry, ok2 := le.(*logrus.Entry); ok2 {
			return entry
		}
	}
	return logrus.NewEntry(logrus.StandardLogger())
}
