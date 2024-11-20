package entity

import (
	"time"
)

type ProductAttribute struct {
	Id        int  `gorm:"type:int;primary_key"`
	ProductId uint `gorm:"null" json:"product_id"`
	Product   Product
	Type      string `gorm:"type:varchar(255)" json:"type"`
	Name      string `gorm:"type:varchar(255)" json:"name"`
	Price     uint   `json:"price"`
	Discount  uint   `json:"discount"`
	Promo     uint   `json:"promo"`
	Image     string `gorm:"type:varchar(255)" json:"image"`
	IsActive  bool   `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
