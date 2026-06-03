// @title           Subscriptions API
// @description     CRUDL API for subscription records.
// @BasePath        /v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            x-api-key
package main

import (
	"em-test-go/src/middlewares"
	"em-test-go/src/models"
	"em-test-go/src/routes"

	_ "em-test-go/docs"

	"github.com/gin-gonic/gin"
)

func main() {
	models.OpenDatabaseConnection()
	models.AutoMigrateModels()

	router := gin.New()
	router.Use(gin.Logger())
	middlewares.RegisterMiddlewares(router)
	routes.SetupRoutes(router)
	router.Run(":8080")

}
