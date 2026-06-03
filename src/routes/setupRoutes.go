package routes

import (
	"em-test-go/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
	}())

	versionRouter := r.Group("/v1")
	versionRouter.GET("/health", controllers.HealthCheck)
	versionRouter.HEAD("/health", controllers.HealthCheck)

	setupSwagger(versionRouter)
	subscriptionsGroupRouter(versionRouter)
	subscriptionCostRouter(versionRouter)
}
