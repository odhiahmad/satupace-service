package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductVariant struct {
	Id               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	Description      *string        `gorm:"type:text" json:"description"`
	BusinessId       uuid.UUID      `gorm:"not null;index:idx_business_sku,unique"`
	SKU              *string        `gorm:"index:idx_business_sku,unique"`
	Business         *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProductId        *uuid.UUID     `gorm:"not null;index" json:"product_id"`
	Product          *Product       `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	BasePrice        *float64       `json:"base_price"`
	SellPrice        *float64       `json:"sell_price"`
	MinimumSales     *int           `json:"minimum_sales"`
	IgnoreStockCheck *bool          `json:"ignore_stock_check"`
	Stock            *int           `json:"stock"`
	TrackStock       *bool          `json:"track_stock"`
	IsAvailable      *bool          `json:"is_available"`
	IsActive         *bool          `json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
