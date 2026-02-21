package entity

import (
	"time"

	"github.com/google/uuid"
)

type RunGroupMember struct {
	Id       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GroupId  uuid.UUID `gorm:"not null;index"`
	UserId   uuid.UUID `gorm:"not null;index"`
	Status   string    `gorm:"type:varchar(50)"` // pending, joined, left
	JoinedAt time.Time
}
