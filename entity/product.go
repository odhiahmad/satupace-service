package entity

import (
	"strconv"
	"time"
)

func (p Product) GetID() string {
	return strconv.Itoa(p.Id)
}

func (p Product) GetCreatedAt() time.Time {
	return p.CreatedAt
}

type Product struct {
	Id                int              `gorm:"primaryKey" json:"id"`
	BusinessId        int              `gorm:"not null" json:"business_id"`
	Business          *Business        `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProductCategoryId int              `gorm:"not null" json:"product_category_id"`
	ProductCategory   *ProductCategory `gorm:"foreignKey:ProductCategoryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	HasVariant        bool             `gorm:"default:false" json:"has_variant"`
	Variants          []ProductVariant `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"variants"`
	Brand             *string          `gorm:"type:varchar(255)" json:"brand"`
	Name              string           `gorm:"type:varchar(255);not null" json:"name"`
	Description       *string          `gorm:"type:text" json:"description"`
	Image             *string          `gorm:"type:text" json:"image"`
	BasePrice         *float64         `json:"base_price"`
	SKU               *string          `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	Stock             *int             `gorm:"default:0" json:"stock"`
	TrackStock        bool             `gorm:"default:false" json:"track_stock"`
	MinimumSales      *int             `gorm:"default:1" json:"minimum_sales"`
	DiscountId        *int             `gorm:"index" json:"discount_id"`
	Discount          *Discount        `gorm:"foreignKey:DiscountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"discount"`
	TaxId             *int             `gorm:"index" json:"tax_id"`
	Tax               *Tax             `gorm:"foreignKey:TaxId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax"`
	UnitId            *int             `gorm:"index" json:"unit_id"`
	Unit              *ProductUnit     `gorm:"foreignKey:UnitId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"unit"`
	ProductPromos     []ProductPromo   `gorm:"foreignKey:ProductId" json:"product_promos"`
	IsAvailable       bool             `gorm:"default:true" json:"is_available"`
	IsActive          bool             `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}
