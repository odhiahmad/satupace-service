package entity

import (
	"time"

	"github.com/google/uuid"
)

type DirectChatMessage struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MatchId   uuid.UUID `gorm:"type:uuid;not null;index"`
	SenderId  uuid.UUID `gorm:"type:uuid;not null;index"`
	Message   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}
