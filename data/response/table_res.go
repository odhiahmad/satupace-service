package response

import (
	"time"

	"github.com/google/uuid"
)

type TableResponse struct {
	Id         uuid.UUID `json:"id"`
	BusinessId uuid.UUID `json:"business_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
