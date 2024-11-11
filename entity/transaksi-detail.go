package entity

import (
	"time"
)

type TransaksiDetail struct {
	Id           int    `gorm:"type:int;primary_key"`
	TransaksiID  int    `gorm:"null"`
	MenuID       int    `gorm:"null"`
	AttributeID  int    `gorm:"null"`
	UkuranMenuID int    `gorm:"null"`
	Total        string `gorm:"type:varchar(255)" json:"total"`
	Discount     string `gorm:"type:varchar(255)" json:"discount"`
	Promo        string `gorm:"type:varchar(255)" json:"promo"`
	Rating       string `gorm:"type:varchar(255)" json:"rating"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *TransaksiDetail) Prepare() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
