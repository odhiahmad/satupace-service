package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserPhoto struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;index"`
	Url       string    `gorm:"type:varchar(255);not null"`
	Type      string    `gorm:"type:varchar(50)"` // profile, run, verification
	IsPrimary bool
	CreatedAt time.Time
}
