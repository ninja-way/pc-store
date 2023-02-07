package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()

		// after request
		latency := time.Since(t)

		log.Infof("[%d] %s %s | %s",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.RequestURI,
			latency,
		)
	}
}
