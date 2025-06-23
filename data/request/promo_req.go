package request

import "time"

type PromoCreate struct {
	BusinessId         int       `json:"business_id" validate:"required"`
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description"`
	Type               string    `json:"type" validate:"required,oneof=percentage fixed"`
	Amount             float64   `json:"amount" validate:"required"`
	RequiredProductIds []int     `json:"required_product_ids"` // untuk kebutuhan business logic, tidak disimpan di DB
	MinQuantity        int       `json:"min_quantity"`
	ProductIds         []int     `json:"product_ids"` // untuk menyimpan relasi ke ProductPromo
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required"`
	IsGlobal           bool      `json:"is_global"`
	IsActive           bool      `json:"is_active"`
}

type PromoUpdate struct {
	Id                 int       `json:"id" validate:"required"`
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description"`
	Type               string    `json:"type" validate:"required,oneof=percentage fixed"`
	Amount             float64   `json:"amount" validate:"required"`
	RequiredProductIds []int     `json:"required_product_ids"`
	MinQuantity        int       `json:"min_quantity"`
	ProductIds         []int     `json:"product_ids"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required"`
	IsGlobal           bool      `json:"is_global"`
	IsActive           bool      `json:"is_active"`
}
