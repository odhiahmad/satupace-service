package response

import "time"

type AuthResponse struct {
	Id          int                     `json:"id"`
	Email       string                  `json:"email"`
	PhoneNumber *string                 `json:"phone_number,omitempty"`
	Token       string                  `json:"token"`
	IsVerified  bool                    `json:"is_verified"`
	IsActive    bool                    `json:"is_active"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	Role        RoleResponse            `json:"role"`
	Business    BusinessResponse        `json:"business"`
	Branch      *BusinessBranchResponse `json:"branch,omitempty"`
	Memberships []MembershipResponse    `json:"memberships,omitempty"` // ✅ Tambahkan ini
}
