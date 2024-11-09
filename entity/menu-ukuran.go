package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MenuUkuran struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	MenuID    uuid.UUID `gorm:"null"`
	Size      string    `gorm:"type:varchar(255)" json:"size"`
	Price     string    `gorm:"type:varchar(255)" json:"price"`
	Stok      string    `gorm:"type:varchar(255)" json:"stok"`
	IsActive  bool      `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *MenuUkuran) Prepare() error {
	u.ID = uuid.NewV4()
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
