package entity

import (
	"time"
)

type Transaction struct {
	Id              int `gorm:"primaryKey" json:"id"`
	BusinessId      int
	Business        Business `gorm:"foreignKey:BusinessId"`
	CustomerId      *int
	Customer        Customer `gorm:"foreignKey:CustomerId"`
	PaymentMethodId *int
	PaymentMethod   PaymentMethod     `gorm:"foreignKey:PaymentMethodId"`
	BillNumber      string            `gorm:"uniqueIndex" json:"bill_number"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionId" json:"items"`
	Total           float64           `json:"total"`
	Discount        float64           `json:"discount"`
	Promo           float64           `json:"promo"`
	Status          string            `gorm:"type:varchar(255)" json:"status"`
	Rating          *float64          `json:"rating"`
	Notes           *string           `gorm:"type:varchar(255)" json:"notes"`
	AmountReceived  *float64          `json:"amount_received"` // Jumlah uang yang diterima dari pelanggan
	Change          *float64          `json:"change"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
