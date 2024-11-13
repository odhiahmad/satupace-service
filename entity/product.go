package entity

import (
	"time"
)

type Product struct {
	Id          int    `gorm:"type:int;primary_key"`
	CustomerId  int    `gorm:"null" json:"customer_id"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Image       string `gorm:"type:varchar(255)" json:"image"`
	Stok        uint   `json:"stok"`
	IsAvailable bool   `gorm:"not null; column:is_available"`
	IsActive    bool   `gorm:"not null; column:is_active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *Product) Prepare() error {
	u.IsActive = true
	return nil
}
