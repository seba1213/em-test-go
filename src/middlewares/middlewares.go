package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Register middleware on the base router
func RegisterMiddlewares(router *gin.Engine) {
	router.Use(AuthMiddleware())
}
