package entity

import (
	"time"

	"github.com/google/uuid"
)

type SafetyLog struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;index"`
	MatchId   uuid.UUID `gorm:"type:uuid;not null;index"`
	Status    string    `gorm:"type:varchar(50);not null"` // reported, blocked
	Reason    string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}
