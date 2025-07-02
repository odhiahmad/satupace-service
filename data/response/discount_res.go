package response

import "time"

type DiscountResponse struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Amount       float64   `json:"amount"`
	IsPercentage bool      `json:"is_percentage"` // true = persen, false = nominal
	IsGlobal     bool      `json:"is_global"`
	IsMultiple   bool      `json:"is_multiple"` // true = diskon berlaku kelipatan
	IsActive     bool      `json:"is_active"`
	StartAt      time.Time `json:"start_at"` // ⛔ Ini tipe time.Time
	EndAt        time.Time `json:"end_at"`
}
