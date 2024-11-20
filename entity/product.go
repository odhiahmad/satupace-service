package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	BusinessID           uint
	Business             Business `gorm:"foreignKey:BusinessID"`
	ProductCategoryID    int
	ProductCategory      ProductCategory `gorm:"foreignKey:ProductCategoryID"`
	ProductSubCategoryID int
	ProductSubCategory   ProductSubCategory `gorm:"foreignKey:ProductSubCategoryID"`
	ProductUnitID        int
	ProductUnit          ProductUnit `gorm:"foreignKey:ProductUnitID"`
	Name                 string      `gorm:"type:varchar(255)" json:"name"`
	Image                string      `gorm:"type:varchar(255)" json:"image"`
	BasePrice            uint        `json:"base_price"`
	Discount             string      `gorm:"type:varchar(255)" json:"discount"`
	Promo                string      `gorm:"type:varchar(255)" json:"promo"`
	IsAvailable          bool        `gorm:"not null; column:is_available"`
	IsActive             bool        `gorm:"not null; column:is_active"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (u *Product) Prepare() error {
	u.IsActive = true
	return nil
}
