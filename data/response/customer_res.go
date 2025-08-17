package response

import (
	"time"

	"github.com/google/uuid"
)

type CustomerResponse struct {
	Id         uuid.UUID `json:"id"`
	BusinessId uuid.UUID `json:"business_id"`
	Name       string    `json:"name"`
	Phone      *string   `json:"phone,omitempty"`
	Email      *string   `json:"email,omitempty"`
	Address    *string   `json:"address,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
