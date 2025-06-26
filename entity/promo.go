package entity

import "time"

type Promo struct {
	Id               int            `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessId       int            `gorm:"not null" json:"business_id"`
	Business         *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name             string         `gorm:"type:varchar(100);not null" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	Type             string         `gorm:"type:varchar(20);not null" json:"type"` // "percentage", "fixed"
	Amount           float64        `gorm:"not null" json:"amount"`
	RequiredProducts []Product      `gorm:"many2many:promo_required_products;" json:"required_products"`
	MinQuantity      int            `json:"min_quantity"`
	ProductPromos    []ProductPromo `gorm:"foreignKey:PromoId" json:"product_promos"`
	StartDate        time.Time      `gorm:"not null" json:"start_date"`
	EndDate          time.Time      `gorm:"not null" json:"end_date"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
