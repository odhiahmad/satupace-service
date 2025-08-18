package response

import (
	"time"

	"github.com/google/uuid"
)

type AuthResponse struct {
	Id          uuid.UUID           `json:"id"`
	Email       string              `json:"email"`
	PhoneNumber string              `json:"phone_number"`
	Token       string              `json:"token"`
	IsVerified  bool                `json:"is_verified"`
	IsActive    bool                `json:"is_active"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Role        RoleResponse        `json:"role"`
	Business    BusinessResponse    `json:"business"`
	Memberships *MembershipResponse `json:"memberships"`
}
