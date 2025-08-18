package entity

import (
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	Id           uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId   uuid.UUID    `gorm:"not null"`
	TerminalId   uuid.UUID    `gorm:"not null"`
	CashierId    uuid.UUID    `gorm:"not null"`
	Cashier      UserBusiness `gorm:"foreignKey:CashierId"`
	OpenedAt     time.Time    `gorm:"not null"`
	ClosedAt     *time.Time
	OpeningCash  float64 `gorm:"not null"`
	ClosingCash  *float64
	TotalSales   *float64
	TotalRefunds *float64
	Status       string `gorm:"type:varchar(10);not null"` // open / closed
	Notes        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
