package entity

import (
	"time"
)

type ProductVariant struct {
	Id          int      `gorm:"primaryKey"`
	BusinessId  int      `gorm:"not null"`
	Business    Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProductId   int      `gorm:"not null"`
	Product     Product  `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name        string   `gorm:"type:varchar(255)"`
	Image       string   `gorm:"type:varchar(255)"`
	BasePrice   float64
	Discount    float64
	Promo       float64
	Stock       int
	FinalPrice  float64
	SKU         string `gorm:"uniqueIndex"`
	IsAvailable bool   `gorm:"not null"`
	IsActive    bool   `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *ProductVariant) Prepare() error {
	p.IsActive = true
	return nil
}
