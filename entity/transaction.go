package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	BusinessID      uint
	Business        Business `gorm:"foreignKey:BusinessID"`
	CustomerId      uint
	Customer        Customer `gorm:"foreignKey:CustomerId"`
	PaymentMethodId uint
	PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodId"`
	BillNumber      string        `gorm:"type:varchar(255)" json:"bill_number"`
	Total           uint          `json:"total"`
	Discount        string        `gorm:"type:varchar(255)" json:"discount"`
	Promo           uint          ` json:"promo"`
	Status          string        `gorm:"type:varchar(255)" json:"status"`
	Rating          uint          ` json:"rating"`
	Notes           string        `gorm:"type:varchar(255)" json:"notes"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
