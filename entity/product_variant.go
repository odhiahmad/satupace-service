package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	Id               int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	Description      *string        `gorm:"type:text" json:"description,omitempty"`
	BusinessId       *int           `gorm:"not null;index:idx_business_sku,unique"`
	SKU              *string        `gorm:"index:idx_business_sku,unique"`
	Business         *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProductId        *int           `gorm:"not null" json:"product_id"`
	Product          *Product       `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	BasePrice        *float64       `json:"base_price"`
	SellPrice        *float64       `json:"sell_price"`
	MinimumSales     *int           `json:"minimum_sales,omitempty"`
	IgnoreStockCheck *bool          `gorm:"default:false" json:"ignore_stock_check"`
	Stock            int            `json:"stock"`
	TrackStock       bool           `gorm:"default:false" json:"track_stock"`
	IsAvailable      *bool          `gorm:"default:true" json:"is_available"`
	IsActive         *bool          `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
