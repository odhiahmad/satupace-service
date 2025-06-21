package entity

import (
	"time"
)

type Product struct {
	Id                int              `gorm:"primaryKey"`
	BusinessId        int              `gorm:"not null"`
	Business          Business         `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProductCategoryId int              `gorm:"not null"`
	ProductCategory   ProductCategory  `gorm:"foreignKey:ProductCategoryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HasVariant        bool             `gorm:"default:false"`
	Variants          []ProductVariant `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name              string           `gorm:"type:varchar(255)"`
	Description       string           `gorm:"type:varchar(255)"`
	Image             string           `gorm:"type:text"`
	BasePrice         float64
	Discount          float64
	Promo             float64
	Stock             int
	FinalPrice        float64
	SKU               string `gorm:"uniqueIndex"`
	IsAvailable       bool   `gorm:"not null"`
	IsActive          bool   `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p *Product) Prepare() error {
	p.IsActive = true
	p.IsAvailable = true
	return nil
}
