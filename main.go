// @title           Subscriptions API
// @description     CRUDL API for subscription records.
// @BasePath        /v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            x-api-key
package main

import (
	"log/slog"

	"em-test-go/src/logging"
	"em-test-go/src/middlewares"
	"em-test-go/src/models"
	"em-test-go/src/routes"

	_ "em-test-go/docs"

	"github.com/gin-gonic/gin"
)

func main() {
	logging.Init()

	models.OpenDatabaseConnection()
	models.AutoMigrateModels()

	router := gin.New()
	router.Use(gin.Recovery())
	middlewares.RegisterMiddlewares(router)
	routes.SetupRoutes(router)
	slog.Info("server starting", "addr", ":8080")
	router.Run(":8080")

}
