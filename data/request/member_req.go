package request

import "github.com/google/uuid"

type MembershipRequest struct {
	UserId uuid.UUID `json:"user_id" validate:"required"`
	Type   string    `json:"type" validate:"required,oneof=monthly yearly weekly"` // hanya boleh 'monthly' atau 'yearly'
}
