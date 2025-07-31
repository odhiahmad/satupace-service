package request

import (
	"github.com/google/uuid"
)

type BrandRequest struct {
	BusinessId uuid.UUID `json:"business_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
}
