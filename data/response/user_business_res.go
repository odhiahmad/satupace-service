package response

import "time"

type UserBusinessResponse struct {
	Id          int                 `json:"id"`
	Email       *string             `json:"email"`
	PhoneNumber string              `json:"phone_number"`
	IsVerified  bool                `json:"is_verified"`
	IsActive    bool                `json:"is_active"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Token       string              `json:"token"`
	Role        *RoleResponse       `json:"role"`
	Business    *BusinessResponse   `json:"business"`
	Membership  *MembershipResponse `json:"membership"`
}
