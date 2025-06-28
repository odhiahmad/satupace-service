package entity

import "time"

type ProductVariant struct {
	Id            int            `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	BusinessId    int            `gorm:"not null" json:"business_id"`
	Business      *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProductId     int            `gorm:"not null" json:"product_id"`
	Product       *Product       `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Image         *string        `gorm:"type:varchar(255)" json:"image"`
	BasePrice     *float64       `json:"base_price"`
	Stock         int            `json:"stock"`
	TrackStock    *bool          `gorm:"default:false" json:"track_stock"`
	SKU           string         `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	DiscountId    *int           `gorm:"index" json:"discount_id"`
	Discount      *Discount      `gorm:"foreignKey:DiscountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"discount"`
	TaxId         *int           `gorm:"index" json:"tax_id"`
	Tax           *Tax           `gorm:"foreignKey:TaxId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax"`
	UnitId        *int           `gorm:"index" json:"unit_id"`
	Unit          *Unit          `gorm:"foreignKey:UnitId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"unit"`
	ProductPromos []ProductPromo `gorm:"foreignKey:ProductVariantId" json:"product_promos"`
	IsAvailable   bool           `gorm:"default:true" json:"is_available"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
