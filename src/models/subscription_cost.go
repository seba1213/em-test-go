package models

import (
	"github.com/google/uuid"
)

// SubscriptionCostResult holds the total cost for a requested period and applied filters.
type SubscriptionCostResult struct {
	Total       int     `json:"total"`
	PeriodStart string  `json:"period_start"`
	PeriodEnd   string  `json:"period_end"`
	UserID      *string `json:"user_id,omitempty"`
	ServiceName string  `json:"service_name,omitempty"`
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
		return nil, err
	}

	total := 0
	for _, sub := range subs {
		months, err := monthsBetween(sub.StartDate, sub.EndDate, params.PeriodStart, params.PeriodEnd)
		if err != nil {
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
