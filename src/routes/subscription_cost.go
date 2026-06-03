package routes

import (
	"em-test-go/src/controllers"

	"github.com/gin-gonic/gin"
)

func subscriptionCostRouter(baseRouter *gin.RouterGroup) {
	baseRouter.GET("/subscription-cost", controllers.GetSubscriptionsTotalCost)
}
