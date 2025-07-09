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
	Token       string              `json:"token,omitempty"` // optional
	Role        *RoleResponse       `json:"role,omitempty"`  // jika ingin tampilkan nama role
	Business    *BusinessResponse   `json:"business,omitempty"`
	Membership  *MembershipResponse `json:"membership,omitempty"`
}
