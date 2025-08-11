package response

import (
	"time"

	"github.com/google/uuid"
)

type TerminalResponse struct {
	Id         uuid.UUID `json:"id"`
	BusinessId uuid.UUID `json:"business_id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
