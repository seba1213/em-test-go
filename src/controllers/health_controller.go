package controllers

import (
	"net/http"

	"em-test-go/src/buildinfo"
	"em-test-go/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetDatabase - returns the database instance
func GetDatabase() *gorm.DB {
	return models.Database
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns service and database health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthSuccessResponse
// @Failure      503  {object}  HealthErrorResponse
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := models.Database.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Database ping failed",
			"error":   err.Error(),
		})
		return
	}

	// Return healthy status
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Service is running",
		"version": buildinfo.Version,
	})
}
