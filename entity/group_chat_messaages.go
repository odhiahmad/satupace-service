package entity

import (
	"time"

	"github.com/google/uuid"
)

type GroupChatMessage struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GroupId   uuid.UUID `gorm:"type:uuid;not null;index"`
	SenderId  uuid.UUID `gorm:"type:uuid;not null;index"`
	Message   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}
