package request

import "github.com/google/uuid"

type UnitRequest struct {
	BusinessId uuid.UUID `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Alias      string    `json:"alias"`
}
