package response

import "github.com/google/uuid"

type OrderTypeResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
