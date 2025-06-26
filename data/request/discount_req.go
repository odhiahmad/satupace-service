package request

import "time"

type DiscountCreate struct {
	BusinessId int       `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Amount     float64   `json:"amount" validate:"required"`
	StartAt    time.Time `json:"start_at" validate:"required"`
	EndAt      time.Time `json:"end_at" validate:"required"`
	IsGlobal   bool      `json:"is_global"`
}

type DiscountUpdate struct {
	BusinessId int       `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Amount     float64   `json:"amount" validate:"required"`
	StartAt    time.Time `json:"start_at" validate:"required"`
	EndAt      time.Time `json:"end_at" validate:"required"`
}

type DiscountIsGlobal struct {
	IsGlobal bool `json:"is_global"`
}

type DiscountIsActive struct {
	IsActive bool `json:"is_active"`
}
