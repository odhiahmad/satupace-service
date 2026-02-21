package request

type CreateUserRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" validate:"required,email"`
	PhoneNumber string  `json:"phone_number" validate:"required"`
	Gender      *string `json:"gender"`
	Password    string  `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Name         *string `json:"name"`
	Gender       *string `json:"gender"`
	PendingEmail *string `json:"pending_email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type VerifyPhoneRequest struct {
	Token string `json:"token" validate:"required"`
}
