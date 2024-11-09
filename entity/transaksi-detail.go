package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TransaksiDetail struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	TransaksiID  uuid.UUID `gorm:"null"`
	MenuID       uuid.UUID `gorm:"null"`
	AttributeID  uuid.UUID `gorm:"null"`
	UkuranMenuID uuid.UUID `gorm:"null"`
	Total        string    `gorm:"type:varchar(255)" json:"total"`
	Discount     string    `gorm:"type:varchar(255)" json:"discount"`
	Promo        string    `gorm:"type:varchar(255)" json:"promo"`
	Rating       string    `gorm:"type:varchar(255)" json:"rating"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *TransaksiDetail) Prepare() error {
	u.ID = uuid.NewV4()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
