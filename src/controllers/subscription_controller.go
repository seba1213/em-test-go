package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"em-test-go/src/logging"
	"em-test-go/src/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func parseSubscriptionID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		logging.WithContext(c).Warn("invalid subscription id", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid subscription id",
		})
		return 0, false
	}
	return uint(id), true
}

func subscriptionNotFound(c *gin.Context) {
	logging.WithContext(c).Warn("subscription not found", "id", c.Param("id"))
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
		logging.WithContext(c).Warn("create subscription: invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	if input.Price < 0 {
		logging.WithContext(c).Warn("create subscription: invalid price", "price", input.Price)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "price must be greater than or equal to 0",
		})
		return
	}
	if !subscriptionDatePattern.MatchString(input.StartDate) {
		logging.WithContext(c).Warn("create subscription: invalid start_date", "start_date", input.StartDate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "start_date must be in MM-YYYY format",
		})
		return
	}
	if input.EndDate != "" && !subscriptionDatePattern.MatchString(input.EndDate) {
		logging.WithContext(c).Warn("create subscription: invalid end_date", "end_date", input.EndDate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "end_date must be in MM-YYYY format",
		})
		return
	}
	if input.EndDate != "" && Base.isEndDateBeforeStartDate(input.StartDate, input.EndDate) {
		logging.WithContext(c).Warn("create subscription: end_date before start_date",
			"start_date", input.StartDate, "end_date", input.EndDate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "end_date must not be before start_date",
		})
		return
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		logging.WithContext(c).Warn("create subscription: invalid user_id", "user_id", input.UserID)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user_id",
		})
		return
	}

	saved, err := input.ToSubscription(userID).Save()
	if err != nil {
		logging.WithContext(c).Error("create subscription: database error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	logging.WithContext(c).Info("subscription created", "subscription_id", saved.ID, "user_id", saved.UserID)
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
		logging.WithContext(c).Error("get subscription: database error", "subscription_id", id, "error", err)
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
			logging.WithContext(c).Warn("list subscriptions: invalid user_id", "user_id", userIDStr)
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
		logging.WithContext(c).Error("list subscriptions: database error", "error", err)
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
		logging.WithContext(c).Warn("update subscription: invalid request body", "subscription_id", id, "error", err)
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
		logging.WithContext(c).Error("update subscription: fetch failed", "subscription_id", id, "error", err)
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
			logging.WithContext(c).Warn("update subscription: invalid price", "subscription_id", id, "price", *input.Price)
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
			logging.WithContext(c).Warn("update subscription: invalid user_id", "subscription_id", id, "user_id", *input.UserID)
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
			logging.WithContext(c).Warn("update subscription: invalid start_date", "subscription_id", id, "start_date", *input.StartDate)
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
			logging.WithContext(c).Warn("update subscription: invalid end_date", "subscription_id", id, "end_date", input.EndDate)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "end_date must be in MM-YYYY format",
			})
			return
		}
		updated.EndDate = input.EndDate
	}
	if updated.EndDate != "" && Base.isEndDateBeforeStartDate(updated.StartDate, updated.EndDate) {
		logging.WithContext(c).Warn("update subscription: end_date before start_date",
			"subscription_id", id, "start_date", updated.StartDate, "end_date", updated.EndDate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "end_date must not be before start_date",
		})
		return
	}

	result, err := updated.UpdateSubscription(id)
	if err != nil {
		logging.WithContext(c).Error("update subscription: database error", "subscription_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	logging.WithContext(c).Info("subscription updated", "subscription_id", result.ID)
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
		logging.WithContext(c).Error("delete subscription: database error", "subscription_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	logging.WithContext(c).Info("subscription deleted", "subscription_id", id)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "subscription deleted",
	})
}
