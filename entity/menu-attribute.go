package entity

import (
	"time"
)

type MenuAttribute struct {
	Id        int    `gorm:"type:int;primary_key"`
	MenuID    int    `gorm:"null"`
	Type      string `gorm:"type:varchar(255)" json:"type"`
	Name      string `gorm:"type:varchar(255)" json:"name"`
	Price     string `gorm:"type:varchar(255)" json:"price"`
	Gambar    string `gorm:"type:varchar(255)" json:"gambar"`
	IsActive  bool   `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *MenuAttribute) Prepare() error {
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
