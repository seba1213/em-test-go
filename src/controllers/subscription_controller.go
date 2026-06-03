package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"em-test-go/src/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func parseSubscriptionID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid subscription id",
		})
		return 0, false
	}
	return uint(id), true
}

func subscriptionNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "failed",
		"message": "subscription not found",
	})
}

// CreateSubscription godoc
// @Summary      Create subscription
// @Description  Creates a subscription record.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      models.SubscriptionToSave  true  "Subscription payload"
// @Success      201   {object}  SubscriptionSuccessResponse
// @Failure      400   {object}  APIBadRequestResponse
// @Failure      401   {object}  APIUnauthorizedResponse
// @Failure      500   {object}  APIInternalErrorResponse
// @Router       /subscriptions/create [post]
func CreateSubscription(c *gin.Context) {
	var input models.SubscriptionToSave
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if input.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "price must be greater than or equal to 0",
		})
		return
	}
	if !subscriptionDatePattern.MatchString(input.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "start_date must be in MM-YYYY format",
		})
		return
	}
	if input.EndDate != "" && !subscriptionDatePattern.MatchString(input.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "end_date must be in MM-YYYY format",
		})
		return
	}
	if input.EndDate != "" && Base.isEndDateBeforeStartDate(input.StartDate, input.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "end_date must not be before start_date",
		})
		return
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user_id",
		})
		return
	}

	saved, err := input.ToSubscription(userID).Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   saved,
	})
}

// GetSubscription godoc
// @Summary      Get subscription
// @Tags         subscriptions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Subscription ID"  default(1)
// @Success      200  {object}  SubscriptionSuccessResponse
// @Failure      400  {object}  APIBadRequestResponse
// @Failure      401  {object}  APIUnauthorizedResponse
// @Failure      404  {object}  APINotFoundResponse
// @Failure      500  {object}  APIInternalErrorResponse
// @Router       /subscriptions/get/{id} [get]
func GetSubscription(c *gin.Context) {
	id, ok := parseSubscriptionID(c)
	if !ok {
		return
	}

	sub, err := models.FetchSubscription(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			subscriptionNotFound(c)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   sub,
	})
}

// ListSubscriptions godoc
// @Summary      List subscriptions
// @Description  Returns a paginated list. Optional filter by user_id.
// @Tags         subscriptions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        user_id  query     string  false  "Filter by user UUID"  default(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Param        page     query     int     false  "Page number"  default(0)
// @Param        size     query     int     false  "Page size"  default(10)
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  APIBadRequestResponse
// @Failure      401      {object}  APIUnauthorizedResponse
// @Failure      500      {object}  APIInternalErrorResponse
// @Router       /subscriptions/list [get]
func ListSubscriptions(c *gin.Context) {
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

	page, err := models.FetchSubscriptions(userID, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   page,
	})
}

// UpdateSubscription godoc
// @Summary      Update subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path      int                          true  "Subscription ID"
// @Param        body  body      models.SubscriptionToUpdate  true  "Fields to update"
// @Success      200   {object}  SubscriptionSuccessResponse
// @Failure      400   {object}  APIBadRequestResponse
// @Failure      401   {object}  APIUnauthorizedResponse
// @Failure      404   {object}  APINotFoundResponse
// @Failure      500   {object}  APIInternalErrorResponse
// @Router       /subscriptions/update/{id} [put]
func UpdateSubscription(c *gin.Context) {
	id, ok := parseSubscriptionID(c)
	if !ok {
		return
	}

	var input models.SubscriptionToUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	existing, err := models.FetchSubscription(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			subscriptionNotFound(c)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	updated := *existing
	if input.ServiceName != nil {
		updated.ServiceName = *input.ServiceName
	}
	if input.Price != nil {
		if *input.Price < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "price must be greater than or equal to 0",
			})
			return
		}
		updated.Price = int(*input.Price)
	}
	if input.UserID != nil {
		userID, err := uuid.Parse(*input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "invalid user_id",
			})
			return
		}
		updated.UserID = userID
	}
	if input.StartDate != nil {
		if !subscriptionDatePattern.MatchString(*input.StartDate) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "start_date must be in MM-YYYY format",
			})
			return
		}
		updated.StartDate = *input.StartDate
	}
	if input.EndDate != "" {
		if !subscriptionDatePattern.MatchString(input.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "end_date must be in MM-YYYY format",
			})
			return
		}
		updated.EndDate = input.EndDate
	}
	if updated.EndDate != "" && Base.isEndDateBeforeStartDate(updated.StartDate, updated.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "end_date must not be before start_date",
		})
		return
	}

	result, err := updated.UpdateSubscription(id)
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

// DeleteSubscription godoc
// @Summary      Delete subscription
// @Tags         subscriptions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Subscription ID"  default(1)
// @Success      200  {object}  SubscriptionDeleteResponse
// @Failure      400  {object}  APIBadRequestResponse
// @Failure      401  {object}  APIUnauthorizedResponse
// @Failure      404  {object}  APINotFoundResponse
// @Failure      500  {object}  APIInternalErrorResponse
// @Router       /subscriptions/delete/{id} [delete]
func DeleteSubscription(c *gin.Context) {
	id, ok := parseSubscriptionID(c)
	if !ok {
		return
	}

	err := models.DeleteSubscription(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			subscriptionNotFound(c)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "subscription deleted",
	})
}
