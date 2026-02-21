package response

import "time"

type UserResponse struct {
	Id          string    `json:"id"`
	Name        *string   `json:"name"`
	Email       *string   `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Gender      *string   `json:"gender"`
	IsVerified  bool      `json:"is_verified"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserDetailResponse struct {
	Id          string    `json:"id"`
	Name        *string   `json:"name"`
	Email       *string   `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Gender      *string   `json:"gender"`
	IsVerified  bool      `json:"is_verified"`
	IsActive    bool      `json:"is_active"`
	Token       string    `json:"token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
