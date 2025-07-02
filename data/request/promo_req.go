package request

import "time"

type PromoCreate struct {
	BusinessId         int       `json:"business_id" validate:"required"`
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description"`
	Type               string    `json:"type" validate:"required"` // e.g. "minimum_spend", "bundle", "buy_x_get_y"
	Amount             float64   `json:"amount" validate:"required"`
	IsPercentage       bool      `json:"is_percentage"` // true jika Amount adalah persen
	MinSpend           *float64  `json:"min_spend"`     // jika Type = "minimum_spend"
	MinQuantity        *int      `json:"min_quantity"`
	RequiredProductIds []int     `json:"required_product_ids"` // produk yang harus dibeli (untuk bundle / buy A+B)
	ProductIds         []int     `json:"product_ids"`          // untuk menyimpan relasi ke ProductPromo
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required"`
	IsActive           bool      `json:"is_active"`
}

type PromoUpdate struct {
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description"`
	Type               string    `json:"type" validate:"required"`
	Amount             float64   `json:"amount" validate:"required"`
	IsPercentage       bool      `json:"is_percentage"`
	MinSpend           *float64  `json:"min_spend"`
	MinQuantity        *int      `json:"min_quantity"`
	RequiredProductIds []int     `json:"required_product_ids"`
	ProductIds         []int     `json:"product_ids"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required"`
	IsActive           bool      `json:"is_active"`
}
