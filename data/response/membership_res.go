package response

import (
	"time"

	"github.com/google/uuid"
)

type MembershipResponse struct {
	Id        uuid.UUID `json:"id"`
	Type      string    `json:"type"` // monthly / yearly
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
}
