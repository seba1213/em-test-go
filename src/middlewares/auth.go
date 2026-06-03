package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func abortUnauthorizedJSON(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  "failed",
		"message": message,
		"data":    nil,
	})
}

func authExcludedPath(path string) bool {
	if path == "/v1/health" {
		return true
	}
	if strings.HasPrefix(path, "/v1/swagger") {
		return true
	}
	return false
}

func AuthMiddleware() gin.HandlerFunc {
	expectedAPIKey := os.Getenv("API_KEY")

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		authRequired := true
		if authExcludedPath(path) || c.Request.Method == "OPTIONS" {
			authRequired = false
		}
		if authRequired {
			apiKey := c.GetHeader("x-api-key")
			if len(apiKey) == 0 {
				abortUnauthorizedJSON(c, "no or empty x-api-key")
				return
			}
			if apiKey != expectedAPIKey {
				abortUnauthorizedJSON(c, "invalid x-api-key")
				return
			}
		}
		c.Next()
	}
}
