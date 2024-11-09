package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MenuAttribute struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	MenuID    uuid.UUID `gorm:"null"`
	Type      string    `gorm:"type:varchar(255)" json:"type"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Price     string    `gorm:"type:varchar(255)" json:"price"`
	Gambar    string    `gorm:"type:varchar(255)" json:"gambar"`
	IsActive  bool      `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *MenuAttribute) Prepare() error {
	u.ID = uuid.NewV4()
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
