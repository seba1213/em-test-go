package models

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

// Subscription stores a user subscription record (user existence is not validated).
type Subscription struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	ServiceName string    `gorm:"not null" json:"service_name"`
	Price       int       `gorm:"not null" json:"price"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	StartDate   string    `gorm:"not null" json:"start_date"`
	EndDate     string    `json:"end_date"`
}

type SubscriptionToSave struct {
	ServiceName string      `json:"service_name" binding:"required"`
	Price       PriceRubles `json:"price" binding:"gte=0" swaggertype:"number" minimum:"0" example:"199.75"`
	UserID      string      `json:"user_id" binding:"required,uuid"`
	StartDate   string      `json:"start_date" binding:"required"`
	EndDate     string      `json:"end_date"`
}

type SubscriptionToUpdate struct {
	ServiceName *string      `json:"service_name"`
	Price       *PriceRubles `json:"price" binding:"omitempty,gte=0" swaggertype:"number" minimum:"0" example:"199.75"`
	UserID      *string      `json:"user_id" binding:"omitempty,uuid"`
	StartDate   *string      `json:"start_date" binding:"omitempty"`
	EndDate     string       `json:"end_date" binding:"omitempty"`
}

func monthYearIndex(date string) (int, error) {
	parts := strings.Split(date, "-")
	if len(parts) != 2 {
		return 0, errors.New("invalid date")
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, errors.New("invalid date")
	}
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("invalid date")
	}
	return year*12 + month, nil
}

func laterMonthIndex(a, b string) (int, error) {
	aIdx, err := monthYearIndex(a)
	if err != nil {
		return 0, err
	}
	bIdx, err := monthYearIndex(b)
	if err != nil {
		return 0, err
	}
	if aIdx >= bIdx {
		return aIdx, nil
	}
	return bIdx, nil
}

func earlierMonthIndex(a, b string) (int, error) {
	aIdx, err := monthYearIndex(a)
	if err != nil {
		return 0, err
	}
	bIdx, err := monthYearIndex(b)
	if err != nil {
		return 0, err
	}
	if aIdx <= bIdx {
		return aIdx, nil
	}
	return bIdx, nil
}

func monthsBetween(subStart, subEnd, periodStart, periodEnd string) (int, error) {
	effectiveStart, err := laterMonthIndex(subStart, periodStart)
	if err != nil {
		return 0, err
	}

	subEndBound := subEnd
	if subEndBound == "" {
		subEndBound = periodEnd
	}

	effectiveEnd, err := earlierMonthIndex(subEndBound, periodEnd)
	if err != nil {
		return 0, err
	}

	if effectiveEnd < effectiveStart {
		return 0, nil
	}
	return effectiveEnd - effectiveStart + 1, nil
}

func (input *SubscriptionToSave) ToSubscription(userID uuid.UUID) *Subscription {
	return &Subscription{
		ServiceName: input.ServiceName,
		Price:       int(input.Price),
		UserID:      userID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}
}

func (sub *Subscription) Save() (*Subscription, error) {
	if err := Database.Create(sub).Error; err != nil {
		return nil, err
	}
	return sub, nil
}

func FetchSubscription(id uint) (*Subscription, error) {
	var sub Subscription
	err := Database.Where("id = ?", id).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func FetchSubscriptions(userID *uuid.UUID, request *http.Request) (paginate.Page, error) {
	pg := paginate.New()
	stmt := Database.Model(&Subscription{})
	if userID != nil {
		stmt = stmt.Where("user_id = ?", *userID)
	}
	page := pg.With(stmt).Request(request).Response(&[]Subscription{})
	return page, nil
}

func (sub *Subscription) UpdateSubscription(id uint) (*Subscription, error) {
	var result Subscription
	err := Database.Model(&Subscription{}).
		Where("id = ?", id).
		Updates(sub).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func DeleteSubscription(id uint) error {
	result := Database.Where("id = ?", id).Delete(&Subscription{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
