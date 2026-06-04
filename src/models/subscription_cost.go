package models

import (
	"log/slog"

	"github.com/google/uuid"
)

// SubscriptionCostResult holds the total cost for a requested period and applied filters.
type SubscriptionCostResult struct {
	Total       int     `json:"total" example:"4800"`
	PeriodStart string  `json:"period_start" example:"01-2024"`
	PeriodEnd   string  `json:"period_end" example:"12-2024"`
	UserID      *string `json:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string  `json:"service_name,omitempty" example:"Yandex Plus"`
}

// SubscriptionCostParams defines period bounds and optional filters for cost calculation.
type SubscriptionCostParams struct {
	PeriodStart string
	PeriodEnd   string
	UserID      *uuid.UUID
	ServiceName string
}

// CalculateSubscriptionsTotalCost sums monthly subscription prices for months overlapping the period.
func CalculateSubscriptionsTotalCost(params SubscriptionCostParams) (*SubscriptionCostResult, error) {
	stmt := Database.Model(&Subscription{})
	if params.UserID != nil {
		stmt = stmt.Where("user_id = ?", *params.UserID)
	}
	if params.ServiceName != "" {
		stmt = stmt.Where("service_name = ?", params.ServiceName)
	}

	var subs []Subscription
	if err := stmt.Find(&subs).Error; err != nil {
		slog.Error("subscription cost query failed", "error", err)
		return nil, err
	}

	total := 0
	for _, sub := range subs {
		months, err := monthsBetween(sub.StartDate, sub.EndDate, params.PeriodStart, params.PeriodEnd)
		if err != nil {
			slog.Warn("subscription cost: skipped subscription with invalid dates",
				"subscription_id", sub.ID, "start_date", sub.StartDate, "end_date", sub.EndDate, "error", err)
			continue
		}
		total += sub.Price * months
	}

	result := &SubscriptionCostResult{
		Total:       total,
		PeriodStart: params.PeriodStart,
		PeriodEnd:   params.PeriodEnd,
		ServiceName: params.ServiceName,
	}
	if params.UserID != nil {
		id := params.UserID.String()
		result.UserID = &id
	}
	return result, nil
}
