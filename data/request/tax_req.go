package request

import "github.com/google/uuid"

type TaxRequest struct {
	BusinessId uuid.UUID `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Amount     float64   `json:"amount" validate:"required"`
	IsGlobal   *bool     `json:"is_global"`
	IsActive   *bool     `json:"is_active"`
}
