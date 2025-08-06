package response

import "github.com/google/uuid"

type UnitResponse struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"` // "Pcs", "Kg", dll
	Alias string    `json:"alias"`
}
