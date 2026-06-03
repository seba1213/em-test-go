package controllers

import "em-test-go/src/models"

// apiErrorData is documented as null in failed API response envelopes.
type apiErrorData struct{}

// APIBadRequestResponse is a 400 failed API response envelope.
type APIBadRequestResponse struct {
	Status  string        `json:"status" example:"failed"`
	Message string        `json:"message" example:"invalid subscription id"`
	Data    *apiErrorData `json:"data" extensions:"x-nullable"`
}

// APIUnauthorizedResponse is a 401 failed API response envelope.
type APIUnauthorizedResponse struct {
	Status  string        `json:"status" example:"failed"`
	Message string        `json:"message" example:"no or empty x-api-key"`
	Data    *apiErrorData `json:"data" extensions:"x-nullable"`
}

// APINotFoundResponse is a 404 failed API response envelope.
type APINotFoundResponse struct {
	Status  string        `json:"status" example:"failed"`
	Message string        `json:"message" example:"subscription not found"`
	Data    *apiErrorData `json:"data" extensions:"x-nullable"`
}

// APIInternalErrorResponse is a 500 failed API response envelope.
type APIInternalErrorResponse struct {
	Status  string        `json:"status" example:"failed"`
	Message string        `json:"message" example:"database error"`
	Data    *apiErrorData `json:"data" extensions:"x-nullable"`
}

// SubscriptionSuccessResponse wraps a single subscription.
type SubscriptionSuccessResponse struct {
	Status string              `json:"status" example:"success"`
	Data   models.Subscription `json:"data"`
}

// SubscriptionDeleteResponse wraps a delete confirmation.
type SubscriptionDeleteResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"subscription deleted"`
}

// SubscriptionCostSuccessResponse wraps total subscription cost for a period.
type SubscriptionCostSuccessResponse struct {
	Status string                        `json:"status" example:"success"`
	Data   models.SubscriptionCostResult `json:"data"`
}

// HealthSuccessResponse is returned when the service is healthy.
type HealthSuccessResponse struct {
	Status  string `json:"status" example:"healthy"`
	Message string `json:"message" example:"Service is running"`
	Version string `json:"version" example:"dev"`
}

// HealthErrorResponse is returned when the service is unhealthy.
type HealthErrorResponse struct {
	Status  string `json:"status" example:"unhealthy"`
	Message string `json:"message" example:"Database ping failed"`
	Error   string `json:"error"`
}
