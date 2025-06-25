package entity

import "time"

type ProductVariant struct {
	Id          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	BusinessId  int       `gorm:"not null" json:"business_id"`
	Business    *Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProductId   int       `gorm:"not null" json:"product_id"`
	Product     *Product  `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Image       *string   `gorm:"type:varchar(255)" json:"image"`
	BasePrice   *float64  `json:"base_price"`
	Discount    *float64  `json:"discount"`
	Promo       *float64  `json:"promo"`
	Stock       int       `json:"stock"`
	FinalPrice  *float64  `json:"final_price"`
	SKU         *string   `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	IsAvailable bool      `gorm:"not null" json:"is_available"`
	IsActive    bool      `gorm:"not null" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
