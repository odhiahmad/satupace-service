package request

type VerifyOTPRequest struct {
	Identifier      string `json:"identifier"`
	Token           string `json:"token"`
	IsResetPassword bool   `json:"is_reset_password"`
}

type RetryOTPRequest struct {
	Identifier string `json:"identifier"`
}

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier"`
}

type ResetPasswordRequest struct {
	Identifier  string `json:"identifier"`
	NewPassword string `json:"password"`
}
