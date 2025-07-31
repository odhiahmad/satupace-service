package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Unit struct {
	Id         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId uuid.UUID      `gorm:"not null;index:idx_business_unit,unique" json:"business_id"`
	Business   *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name       string         `gorm:"type:varchar(50);not null;index:idx_business_unit,unique" json:"name"`
	Alias      string         `gorm:"type:varchar(20)" json:"alias"`
	Multiplier float64        `gorm:"default:1" json:"multiplier"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
