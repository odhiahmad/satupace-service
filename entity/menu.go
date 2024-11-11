package entity

import (
	"time"
)

type Menu struct {
	Id           int    `gorm:"type:int;primary_key"`
	PerusahaanID int    `gorm:"null"`
	Nama         string `gorm:"type:varchar(255)" json:"nama"`
	Gambar       string `gorm:"type:varchar(255)" json:"gambar"`
	Price        string `gorm:"type:varchar(255)" json:"price"`
	Discount     string `gorm:"type:varchar(255)" json:"discount"`
	Promo        string `gorm:"type:varchar(255)" json:"promo"`
	Stok         string `gorm:"type:varchar(255)" json:"stok"`
	IsActive     bool   `gorm:"not null; column:is_active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *Menu) Prepare() error {
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
