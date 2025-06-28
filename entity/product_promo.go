package entity

import "time"

type ProductPromo struct {
	Id               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductId        *int      `json:"product_id,omitempty"`
	ProductVariantId *int      `json:"product_variant_id,omitempty"`
	BusinessId       int       `gorm:"not null" json:"business_id"`
	Business         *Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	PromoId          int       `gorm:"primaryKey" json:"promo_id"`
	Product          Product   `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Promo            *Promo    `gorm:"foreignKey:PromoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MinQuantity      int       `gorm:"column:min_quantity" json:"min_quantity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
