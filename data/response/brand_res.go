package response

import "github.com/google/uuid"

type BrandResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"` // "Pcs", "Kg", dll
}
