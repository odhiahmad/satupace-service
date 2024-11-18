package entity

import (
	"time"
)

type Transaction struct {
	Id              int    `gorm:"type:int;primary_key"`
	BusinessId      int    `gorm:"null" json:"business_id"`
	CustomerId      int    `gorm:"null" json:"customer_id"`
	BillNumber      string `gorm:"type:varchar(255)" json:"bill_number"`
	PaymentMethodId string `gorm:"type:varchar(255)" json:"payment_method_id"`
	Total           uint   `json:"total"`
	Discount        string `gorm:"type:varchar(255)" json:"discount"`
	Promo           uint   ` json:"promo"`
	Status          string `gorm:"type:varchar(255)" json:"status"`
	Rating          uint   ` json:"rating"`
	Notes           string `gorm:"type:varchar(255)" json:"notes"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
