package response

import "time"

type DiscountResponse struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Amount       float64   `json:"amount"`
	IsPercentage *bool     `json:"is_percentage"`
	IsGlobal     *bool     `json:"is_global"`
	IsMultiple   *bool     `json:"is_multiple"`
	IsActive     *bool     `json:"is_active"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}
