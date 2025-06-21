package entity

import (
	"time"
)

type Transaction struct {
	Id              int `gorm:"type:int;primary_key" json:"transaction_id"`
	BusinessId      int
	Business        Business `gorm:"foreignKey:BusinessId"`
	CustomerId      *int
	Customer        Customer `gorm:"foreignKey:CustomerId"`
	PaymentMethodId *int
	PaymentMethod   PaymentMethod     `gorm:"foreignKey:PaymentMethodId"`
	BillNumber      string            `gorm:"type:varchar(255)" json:"bill_number"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionId"`
	Total           int               `json:"total"`
	Discount        string            `gorm:"type:varchar(255)" json:"discount"`
	Promo           int               ` json:"promo"`
	Status          string            `gorm:"type:varchar(255)" json:"status"`
	Rating          int               ` json:"rating"`
	Notes           string            `gorm:"type:varchar(255)" json:"notes"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
