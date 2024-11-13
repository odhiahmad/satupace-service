package entity

import (
	"time"
)

type Customer struct {
	User        User
	Name        string `gorm:"type:varchar(255)" json:"name"`
	PhoneNumber string `gorm:"uniqueIndex;type:varchar(255)" json:"phone_number"`
	IsActive    bool   `gorm:"not null; column:is_active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *Customer) Prepare() error {
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
