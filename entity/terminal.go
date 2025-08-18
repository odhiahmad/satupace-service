package entity

import (
	"time"

	"github.com/google/uuid"
)

type Terminal struct {
	Id         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId uuid.UUID `gorm:"not null"`
	Business   Business  `gorm:"foreignKey:BusinessId"`
	Name       string    `gorm:"type:varchar(100);not null"`
	Location   string    `gorm:"type:varchar(255)"`
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
