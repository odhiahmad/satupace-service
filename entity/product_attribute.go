package entity

import (
	"time"
)

type ProductAttribute struct {
	Id          int     `gorm:"type:int;primary_key"`
	ProductId   int     `gorm:"null" json:"product_id"`
	Name        string  `gorm:"type:varchar(255)" json:"name"`
	Price       float64 `json:"price"`
	Image       string  `gorm:"type:varchar(255)" json:"image"`
	IsActive    bool    `gorm:"not null; column:is_active"`
	IsAvailable bool    `gorm:"not null; column:is_available"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
