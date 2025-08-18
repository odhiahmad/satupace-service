package request

import "github.com/google/uuid"

type TableRequest struct {
	BusinessId uuid.UUID `json:"business_id" validate:"required"`
	Number     string    `json:"number" validate:"required"`
	Status     string    `json:"status" validate:"omitempty,oneof=available occupied reserved"`
}
