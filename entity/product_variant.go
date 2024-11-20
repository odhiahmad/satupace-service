package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	gorm.Model
	ProductId       int     `gorm:"null"`
	Product         Product `gorm:"foreignKey:ProductId"`
	Name            string  `gorm:"type:varchar(255)" json:"name"`
	AdditionalPrice uint    `json:"additional_price"`
	Discount        string  `gorm:"type:varchar(255)" json:"discount"`
	Promo           string  `gorm:"type:varchar(255)" json:"promo"`
	IsAvailable     bool    `gorm:"not null; column:is_active"`
	IsActive        bool    `gorm:"not null; column:is_active"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
