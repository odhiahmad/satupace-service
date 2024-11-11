package entity

import (
	"time"
)

type Transaksi struct {
	Id           int    `gorm:"type:int;primary_key"`
	PerusahaanID int    `gorm:"null"`
	PelangganID  int    `gorm:"null"`
	Total        string `gorm:"type:varchar(255)" json:"total"`
	Discount     string `gorm:"type:varchar(255)" json:"discount"`
	Promo        string `gorm:"type:varchar(255)" json:"promo"`
	Status       string `gorm:"type:varchar(255)" json:"status"`
	Rating       string `gorm:"type:varchar(255)" json:"rating"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *Transaksi) Prepare() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
