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
	CashierId       *int              `json:"user_id"`
	Cashier         *UserBusiness     `gorm:"foreignKey:CashierId;references:Id"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionId" json:"items"`
	FinalPrice      float64           `json:"finalPrice"`
	BasePrice       float64           `json:"basePrice"`
	SellPrice       float64           `json:"sellPrice"`
	Discount        float64           `json:"discount"`
	Promo           float64           `json:"promo"`
	Tax             float64           `json:"tax"`
	Status          string            `gorm:"type:varchar(255)" json:"status"`
	Rating          *float64          `json:"rating"`
	Notes           *string           `gorm:"type:varchar(255)" json:"notes"`
	AmountReceived  *float64          `json:"amount_received"` // Jumlah uang yang diterima dari pelanggan
	Change          *float64          `json:"change"`
	PaidAt          *time.Time        `json:"paid_at"`
	IsRefunded      *bool             `gorm:"default:false" json:"is_refunded"`
	RefundReason    *string           `json:"refund_reason"`
	RefundedAt      *time.Time        `json:"refunded_at"`
	RefundedBy      *int              `json:"refunded_by"`
	IsCanceled      *bool             `gorm:"default:false" json:"is_canceled"`
	CanceledAt      *time.Time        `json:"canceled_at"`
	CanceledBy      *int              `json:"canceled_by"`
	CanceledReason  *string           `json:"canceled_reason"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
