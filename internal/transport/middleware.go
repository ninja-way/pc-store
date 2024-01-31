package transport

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ninja-way/pc-store/internal/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
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

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c)
		if err != nil {
			config.LogDebug("authMiddleware", err)
			c.Status(http.StatusUnauthorized)
			return
		}

		var userService Users
		userID, err := userService.ParseToken(token)
		if err != nil {
			config.LogDebug("authMiddleware", err)
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Set("uid", userID)
		c.Next()
	}
}

func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Values("Authorization")
	if header == nil {
		return "", errors.New("empty auth header")
	}

	authInfo := strings.Split(header[0], " ")
	if len(authInfo) != 2 || authInfo[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(authInfo[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return authInfo[1], nil
}
