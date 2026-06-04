package logging

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func Init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}

func WithContext(c *gin.Context) *slog.Logger {
	if c == nil {
		return slog.Default()
	}
	if rid, exists := c.Get("request_id"); exists {
		if s, ok := rid.(string); ok && s != "" {
			return slog.Default().With("request_id", s)
		}
	}
	return slog.Default()
}
