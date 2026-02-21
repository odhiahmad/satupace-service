package request

type CreateUserRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"required,email"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Gender      *string `json:"gender"`
	Password    string  `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name         *string `json:"name"`
	Gender       *string `json:"gender"`
	PendingEmail *string `json:"pending_email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyPhoneRequest struct {
	Token string `json:"token" binding:"required"`
}
