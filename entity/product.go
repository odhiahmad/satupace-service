package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Id               uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SKU              *string          `gorm:"index:idx_business_sku,unique"`
	BusinessId       uuid.UUID        `gorm:"not null;index:idx_business_sku,unique"`
	Business         *Business        `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CategoryId       *uuid.UUID       `gorm:"index" json:"category_id"`
	Category         *Category        `gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	HasVariant       bool             `json:"has_variant"`
	Variants         []ProductVariant `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"variants"`
	Name             string           `gorm:"type:varchar(255);not null;index" json:"name"`
	Description      *string          `gorm:"type:text" json:"description"`
	Image            *string          `gorm:"type:text" json:"image"`
	BasePrice        *float64         `json:"base_price"`
	SellPrice        *float64         `json:"sell_price"`
	Stock            *int             `json:"stock"`
	TrackStock       *bool            `json:"track_stock"`
	IgnoreStockCheck *bool            `json:"ignore_stock_check"`
	MinimumSales     *int             `json:"minimum_sales"`
	DiscountId       *uuid.UUID       `gorm:"index" json:"discount_id"`
	Discount         *Discount        `gorm:"foreignKey:DiscountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"discount"`
	BrandId          *uuid.UUID       `gorm:"index" json:"brand_id"`
	Brand            *Brand           `gorm:"foreignKey:BrandId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"brand"`
	TaxId            *uuid.UUID       `gorm:"index" json:"tax_id"`
	Tax              *Tax             `gorm:"foreignKey:TaxId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax"`
	UnitId           *uuid.UUID       `gorm:"index" json:"unit_id"`
	Unit             *Unit            `gorm:"foreignKey:UnitId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"unit"`
	IsAvailable      *bool            `json:"is_available"`
	IsActive         *bool            `json:"is_active"`
	IsReady          bool             `gorm:"default:false" json:"is_ready"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        gorm.DeletedAt   `gorm:"index" json:"-"`
}
