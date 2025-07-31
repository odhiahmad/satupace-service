package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	Id              uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId      uuid.UUID         `gorm:"not null;index:idx_business_bill,unique" json:"business_id"`
	Business        *Business         `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CustomerId      *uuid.UUID        `json:"customer_id"`
	Customer        *Customer         `gorm:"foreignKey:CustomerId"`
	PaymentMethodId *int              `json:"payment_method_id"`
	PaymentMethod   PaymentMethod     `gorm:"foreignKey:PaymentMethodId"`
	BillNumber      string            `gorm:"index:idx_business_bill,unique" json:"bill_number"`
	CashierId       *uuid.UUID        `json:"user_id"`
	Cashier         *UserBusiness     `gorm:"foreignKey:CashierId;references:Id"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionId" json:"items"`
	FinalPrice      float64           `json:"final_price"`
	BasePrice       float64           `json:"base_price"`
	SellPrice       float64           `json:"sell_price"`
	Discount        float64           `json:"discount"`
	Promo           float64           `json:"promo"`
	Tax             float64           `json:"tax"`
	Status          string            `gorm:"type:varchar(255);index" json:"status"`
	Rating          *float64          `json:"rating"`
	Notes           *string           `gorm:"type:varchar(255)" json:"notes"`
	AmountReceived  *float64          `json:"amount_received"`
	Change          *float64          `json:"change"`
	PaidAt          *time.Time        `json:"paid_at"`
	IsRefunded      bool              `gorm:"default:false" json:"is_refunded"`
	RefundReason    *string           `json:"refund_reason"`
	RefundedAt      *time.Time        `json:"refunded_at"`
	RefundedBy      *uuid.UUID        `json:"refunded_by"`
	IsCanceled      bool              `gorm:"default:false" json:"is_canceled"`
	CanceledAt      *time.Time        `json:"canceled_at"`
	CanceledBy      *uuid.UUID        `json:"canceled_by"`
	CanceledReason  *string           `json:"canceled_reason"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
}
