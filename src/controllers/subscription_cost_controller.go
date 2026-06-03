package controllers

import (
	"net/http"

	"em-test-go/src/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetSubscriptionsTotalCost godoc
// @Summary      Total subscription cost for a period
// @Description  Sums monthly prices for subscriptions active during the period. Optional filters: user_id, service_name.
// @Tags         subscription-cost
// @Produce      json
// @Security     ApiKeyAuth
// @Param        period_start   query     string  true   "Period start (MM-YYYY)"  default(01-2024)
// @Param        period_end     query     string  true   "Period end (MM-YYYY)"  default(12-2024)
// @Param        user_id        query     string  false  "Filter by user UUID"  default(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Param        service_name   query     string  false  "Filter by subscription service name"  default(Yandex Plus)
// @Success      200            {object}  SubscriptionCostSuccessResponse
// @Failure      400            {object}  APIBadRequestResponse
// @Failure      401            {object}  APIUnauthorizedResponse
// @Failure      500            {object}  APIInternalErrorResponse
// @Router       /subscription-cost [get]
func GetSubscriptionsTotalCost(c *gin.Context) {
	periodStart := c.Query("period_start")
	periodEnd := c.Query("period_end")
	serviceName := c.Query("service_name")

	if periodStart == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "period_start is required",
		})
		return
	}
	if periodEnd == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "period_end is required",
		})
		return
	}
	if !subscriptionDatePattern.MatchString(periodStart) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "period_start must be in MM-YYYY format",
		})
		return
	}
	if !subscriptionDatePattern.MatchString(periodEnd) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "period_end must be in MM-YYYY format",
		})
		return
	}
	if Base.isEndDateBeforeStartDate(periodStart, periodEnd) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "period_end must not be before period_start",
		})
		return
	}

	var userID *uuid.UUID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "invalid user_id query parameter",
			})
			return
		}
		userID = &parsed
	}

	result, err := models.CalculateSubscriptionsTotalCost(models.SubscriptionCostParams{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		UserID:      userID,
		ServiceName: serviceName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}
