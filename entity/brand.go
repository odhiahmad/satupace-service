package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Brand struct {
	Id         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId uuid.UUID      `gorm:"not null;index" json:"business_id"`
	Business   *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name       string         `gorm:"type:varchar(100);not null;index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
