package entity

import "time"

type ProductPromo struct {
	ProductId   int       `gorm:"primaryKey" json:"product_id"`
	BusinessId  int       `gorm:"not null" json:"business_id"`
	Business    *Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	PromoId     int       `gorm:"primaryKey" json:"promo_id"`
	Product     Product   `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Promo       *Promo    `gorm:"foreignKey:PromoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MinQuantity int       `json:"min_quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
