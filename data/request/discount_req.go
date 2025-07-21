package request

import "time"

type DiscountRequest struct {
	BusinessId  int       `json:"business_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount" validate:"required"`
	IsGlobal    *bool     `json:"is_global"`
	IsMultiple  *bool     `json:"is_multiple"`
	StartAt     time.Time `json:"start_at" validate:"required"`
	EndAt       time.Time `json:"end_at" validate:"required"`
	IsActive    *bool     `json:"is_active"`
}

type DiscountIsGlobal struct {
	IsGlobal *bool `json:"is_global"`
}

type DiscountIsActive struct {
	IsActive *bool `json:"is_active"`
}
