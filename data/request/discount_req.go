package request

import "time"

type DiscountCreate struct {
	BusinessId int       `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Type       string    `json:"type" binding:"required,oneof=percent fixed"`
	Amount     float64   `json:"amount" validate:"required"`
	StartAt    time.Time `json:"start_at" validate:"required"` // ✅ Ubah jadi time.Time
	EndAt      time.Time `json:"end_at" validate:"required"`   // ✅ Ubah jadi time.Time
	IsGlobal   bool      `json:"is_global"`
	ProductIds []int     `json:"product_ids,omitempty"`
}

type DiscountUpdate struct {
	Id         int       `json:"id" validate:"required"`
	BusinessId int       `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Type       string    `json:"type" validate:"required,oneof=percent fixed"`
	Amount     float64   `json:"amount" validate:"required"`
	StartAt    time.Time `json:"start_at" validate:"required"`
	EndAt      time.Time `json:"end_at" validate:"required"`
	IsGlobal   bool      `json:"is_global"`
	ProductIds []int     `json:"product_ids,omitempty"`
}
