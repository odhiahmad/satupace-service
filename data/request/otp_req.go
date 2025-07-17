package request

type VerifyOTPRequest struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

type RetryOTPRequest struct {
	Identifier string `json:"identifier"`
}

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier"`
}

type ResetPasswordRequest struct {
	Identifier  string `json:"identifier"`
	OTP         string `json:"otp"`
	NewPassword string `json:"password"`
}
