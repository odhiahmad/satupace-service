package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserBiometric struct {
	Id           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	CredentialId string     `gorm:"type:text;not null;uniqueIndex" json:"credential_id"`
	PublicKey    string     `gorm:"type:text;not null" json:"public_key"`
	DeviceName   string     `gorm:"type:varchar(255)" json:"device_name"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
