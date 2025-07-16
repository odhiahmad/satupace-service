package request

type VerifyOTPRequest struct {
	Via        string `json:"via"`        // "whatsapp", "email", atau "change_email"
	Identifier string `json:"identifier"` // No HP atau Email (tergantung via)
	Token      string `json:"token"`      // Kode OTP
}

type RetryOTPRequest struct {
	Via   string `json:"via"`   // "whatsapp" atau "email"
	Value string `json:"value"` // nomor WA atau alamat email
}

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier"` // email atau no HP
}

type ResetPasswordRequest struct {
	Identifier  string `json:"identifier"`   // email atau no HP
	OTP         string `json:"otp"`          // kode OTP
	NewPassword string `json:"new_password"` // password baru
}
