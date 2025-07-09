package response

import "time"

type MembershipResponse struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"` // monthly / yearly
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
}
