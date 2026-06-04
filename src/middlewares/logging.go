package middlewares

import (
	"strings"
	"time"

	"em-test-go/src/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func loggingExcludedPath(path string) bool {
	return path == "/v1/health"
}

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if loggingExcludedPath(path) {
			c.Next()
			return
		}

		start := time.Now()
		requestID := uuid.NewString()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()

		query := c.Request.URL.RawQuery
		if query != "" {
			query = "?" + query
		}

		logging.WithContext(c).Info("request",
			"method", c.Request.Method,
			"path", path+query,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", strings.TrimSpace(c.Request.UserAgent()),
		)
	}
}
