package entity

import (
	"time"

	"github.com/google/uuid"
)

type DirectMatch struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	User1Id   uuid.UUID `gorm:"not null;index"`
	User2Id   uuid.UUID `gorm:"not null;index"`
	Status    string    `gorm:"type:varchar(50)"` // pending, accepted, rejected
	CreatedAt time.Time
	MatchedAt *time.Time
}
