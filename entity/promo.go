package entity

import "time"

type Promo struct {
	Id               int            `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessId       int            `gorm:"not null" json:"business_id"`
	Business         Business       `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name             string         `gorm:"type:varchar(100);not null" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	Type             string         `gorm:"type:varchar(30);not null" json:"type"` // e.g. minimum_spend, minimum_quantity
	Amount           float64        `gorm:"not null" json:"amount"`
	IsPercentage     bool           `gorm:"not null;default:false" json:"is_percentage"`
	MinSpend         *float64       `json:"min_spend,omitempty"`
	MinQuantity      *int           `json:"min_quantity,omitempty"`
	RequiredProducts []Product      `gorm:"many2many:promo_required_products;joinForeignKey:PromoId;joinReferences:ProductId" json:"required_products,omitempty"`
	ProductPromos    []ProductPromo `gorm:"foreignKey:PromoId" json:"product_promos"`
	StartDate        time.Time      `gorm:"not null" json:"start_date"`
	EndDate          time.Time      `gorm:"not null" json:"end_date"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
