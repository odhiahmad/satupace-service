package entity

import (
	"time"
)

type PerusahaanUser struct {
	Id           int    `gorm:"type:int;primary_key"`
	PerusahaanID int    `gorm:"null"`
	Nama         string `gorm:"type:varchar(255)" json:"nama"`
	Username     string `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
	Password     string `gorm:"->;<-;not null" json:"-"`
	Token        string `gorm:"-" json:"token,omitempty"`
	IsActive     bool   `gorm:"not null; column:is_active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *PerusahaanUser) Prepare() error {
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
