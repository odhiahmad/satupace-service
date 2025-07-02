package response

import "time"

type UserBusinessResponse struct {
	Id          int                     `json:"id"`
	RoleId      int                     `json:"role_id"`
	BusinessId  int                     `json:"business_id"`
	BranchId    *int                    `json:"branch_id,omitempty"`
	Email       string                  `json:"email"`
	PhoneNumber *string                 `json:"phone_number,omitempty"`
	IsVerified  bool                    `json:"is_verified"`
	IsActive    bool                    `json:"is_active"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	Token       string                  `json:"token,omitempty"` // optional
	Role        *RoleResponse           `json:"role,omitempty"`  // jika ingin tampilkan nama role
	Business    *BusinessResponse       `json:"business,omitempty"`
	Branch      *BusinessBranchResponse `json:"branch,omitempty"`
	Memberships []MembershipResponse    `json:"memberships,omitempty"`
}
