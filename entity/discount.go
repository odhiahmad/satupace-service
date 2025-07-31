package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Discount struct {
	Id           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId   uuid.UUID      `gorm:"not null;index" json:"business_id"`
	Business     *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	IsPercentage *bool          `gorm:"not null;default:false" json:"is_percentage"` // true = amount sebagai persen
	Amount       float64        `gorm:"not null" json:"amount"`                      // nilai diskon
	IsGlobal     *bool          `gorm:"not null;default:false" json:"is_global"`     // true = berlaku untuk semua produk
	IsMultiple   *bool          `gorm:"not null;default:false" json:"is_multiple"`   // true = berlaku kelipatan (misal beli 2x, diskon 2x)
	StartAt      time.Time      `json:"start_at"`
	EndAt        time.Time      `json:"end_at"`
	IsActive     *bool          `gorm:"not null;default:false" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
