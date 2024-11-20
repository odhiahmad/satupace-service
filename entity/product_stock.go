package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductStock struct {
	gorm.Model
	ProductSizeID    uint
	ProductSize      ProductSize `gorm:"foreignKey:ProductSizeID"`
	ProductVariantID uint
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID"`
	Quantity         uint           `json:"quantity"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
