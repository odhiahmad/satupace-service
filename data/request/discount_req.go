package request

import "time"

type DiscountRequest struct {
	BusinessId   int       `json:"business_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description"`   // baru
	IsPercentage bool      `json:"is_percentage"` // baru
	Amount       float64   `json:"amount" validate:"required"`
	IsGlobal     bool      `json:"is_global"`
	IsMultiple   bool      `json:"is_multiple"` // baru
	StartAt      time.Time `json:"start_at" validate:"required"`
	EndAt        time.Time `json:"end_at" validate:"required"`
	IsActive     bool      `json:"is_active"` // baru
}

type DiscountIsGlobal struct {
	IsGlobal bool `json:"is_global"`
}

type DiscountIsActive struct {
	IsActive bool `json:"is_active"`
}
