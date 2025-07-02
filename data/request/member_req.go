package request

type MembershipRequest struct {
	UserId int    `json:"user_id" binding:"required"`
	Type   string `json:"type" binding:"required,oneof=monthly yearly weekly"` // hanya boleh 'monthly' atau 'yearly'
}
