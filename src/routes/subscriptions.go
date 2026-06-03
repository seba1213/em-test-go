package routes

import (
	"em-test-go/src/controllers"

	"github.com/gin-gonic/gin"
)

func subscriptionsGroupRouter(baseRouter *gin.RouterGroup) {
	subscriptions := baseRouter.Group("/subscriptions")

	subscriptions.POST("/create", controllers.CreateSubscription)
	subscriptions.GET("/get/:id", controllers.GetSubscription)
	subscriptions.GET("/list", controllers.ListSubscriptions)
	subscriptions.PUT("/update/:id", controllers.UpdateSubscription)
	subscriptions.DELETE("/delete/:id", controllers.DeleteSubscription)
}
