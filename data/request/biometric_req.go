package request

type BiometricRegisterStartRequest struct {
	DeviceName string `json:"device_name" binding:"required"`
}

type BiometricRegisterFinishRequest struct {
	CredentialId string `json:"credential_id" binding:"required"`
	PublicKey    string `json:"public_key" binding:"required"`
	DeviceName   string `json:"device_name" binding:"required"`
	Challenge    string `json:"challenge" binding:"required"`
	Signature    string `json:"signature" binding:"required"`
}

type BiometricLoginStartRequest struct {
	Identifier string `json:"identifier" binding:"required"`
}

type BiometricLoginFinishRequest struct {
	CredentialId string `json:"credential_id" binding:"required"`
	Challenge    string `json:"challenge" binding:"required"`
	Signature    string `json:"signature" binding:"required"`
}
