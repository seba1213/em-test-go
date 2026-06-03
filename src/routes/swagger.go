package routes

import (
	"net/http"
	"strings"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

const swaggerUIPath = "/v1/swagger/index.html"

func setupSwagger(v1 *gin.RouterGroup) {
	swaggerHandler := ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/v1/swagger/doc.json"),
	)

	v1.GET("/swagger/*any", func(c *gin.Context) {
		any := strings.TrimPrefix(c.Param("any"), "/")
		if any == "" {
			c.Redirect(http.StatusMovedPermanently, swaggerUIPath)
			return
		}
		swaggerHandler(c)
	})
}
