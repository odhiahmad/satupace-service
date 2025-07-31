package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tax struct {
	Id           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId   uuid.UUID      `gorm:"not null;index:idx_business_tax,unique" json:"business_id"`
	Business     *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name         string         `gorm:"type:varchar(100);not null;index:idx_business_tax,unique" json:"name"`
	IsPercentage *bool          `json:"is_percentage"`
	IsGlobal     *bool          `json:"is_global"`
	Amount       float64        `gorm:"not null" json:"amount"`
	IsActive     *bool          `json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
