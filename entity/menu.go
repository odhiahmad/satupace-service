package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Menu struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	PerusahaanID uuid.UUID `gorm:"null"`
	Nama         string    `gorm:"type:varchar(255)" json:"nama"`
	Gambar       string    `gorm:"type:varchar(255)" json:"gambar"`
	Price        string    `gorm:"type:varchar(255)" json:"price"`
	Discount     string    `gorm:"type:varchar(255)" json:"discount"`
	Promo        string    `gorm:"type:varchar(255)" json:"promo"`
	Stok         string    `gorm:"type:varchar(255)" json:"stok"`
	IsActive     bool      `gorm:"not null; column:is_active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *Menu) Prepare() error {
	u.ID = uuid.NewV4()
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
