package entity

import (
	"time"
)

type ProductSize struct {
	Id          int    `gorm:"type:int;primary_key"`
	ProductId   int    `gorm:"null"`
	Size        string `gorm:"type:varchar(255)" json:"size"`
	Price       uint   `json:"price"`
	Stok        uint   `json:"stok"`
	IsAvailable bool   `gorm:"not null; column:is_active"`
	IsActive    bool   `gorm:"not null; column:is_active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *ProductSize) Prepare() error {
	u.IsActive = true
	return nil
}
