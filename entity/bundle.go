package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bundle struct {
	Id          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId  uuid.UUID      `gorm:"not null;index" json:"business_id"`
	Business    *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name        string         `gorm:"type:varchar(255)"`
	Description *string        `gorm:"type:varchar(255)"`
	Image       *string        `gorm:"type:text"`
	BasePrice   *float64       `json:"base_price"`
	SellPrice   *float64       `json:"sell_price"`
	Stock       *int           `json:"stock"`
	Items       []BundleItem   `gorm:"foreignKey:BundleId"`
	TaxId       *uuid.UUID     `gorm:"index" json:"tax_id"`
	Tax         *Tax           `gorm:"foreignKey:TaxId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax"`
	IsAvailable bool           `gorm:"not null;default:true" json:"is_available,omitempty"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
