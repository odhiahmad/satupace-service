package request

import (
	"github.com/google/uuid"
)

type CustomerRequest struct {
	BusinessId uuid.UUID `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required,min=2,max=100"`
	Phone      *string   `json:"phone" validate:"omitempty,max=20"`
	Email      *string   `json:"email" validate:"omitempty,email,max=100"`
	Address    *string   `json:"address" validate:"omitempty"`
}
