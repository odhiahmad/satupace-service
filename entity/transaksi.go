package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Transaksi struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	PerusahaanID uuid.UUID `gorm:"null"`
	PelangganID  uuid.UUID `gorm:"null"`
	Total        string    `gorm:"type:varchar(255)" json:"total"`
	Discount     string    `gorm:"type:varchar(255)" json:"discount"`
	Promo        string    `gorm:"type:varchar(255)" json:"promo"`
	Status       string    `gorm:"type:varchar(255)" json:"status"`
	Rating       string    `gorm:"type:varchar(255)" json:"rating"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *Transaksi) Prepare() error {
	u.ID = uuid.NewV4()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
