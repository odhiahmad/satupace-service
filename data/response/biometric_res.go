package response

import "time"

type BiometricChallengeResponse struct {
	Challenge string `json:"challenge"`
}

type BiometricCredentialResponse struct {
	Id           string     `json:"id"`
	UserId       string     `json:"user_id"`
	CredentialId string     `json:"credential_id"`
	DeviceName   string     `json:"device_name"`
	IsActive     bool       `json:"is_active"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
