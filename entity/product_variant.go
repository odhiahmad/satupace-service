package entity

import "time"

type ProductVariant struct {
	Id          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	BusinessId  int       `gorm:"not null" json:"business_id"`
	Business    *Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProductId   int       `gorm:"not null" json:"product_id"`
	Product     *Product  `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	BasePrice   *float64  `json:"base_price"`
	SellPrice   *float64  `json:"sale_price"`
	Stock       int       `json:"stock"`
	TrackStock  bool      `gorm:"default:false" json:"track_stock"`
	SKU         string    `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	IsAvailable bool      `gorm:"default:true" json:"is_available"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
