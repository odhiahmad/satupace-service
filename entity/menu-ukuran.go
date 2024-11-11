package entity

import (
	"time"
)

type MenuUkuran struct {
	Id        int    `gorm:"type:int;primary_key"`
	MenuID    int    `gorm:"null"`
	Size      string `gorm:"type:varchar(255)" json:"size"`
	Price     string `gorm:"type:varchar(255)" json:"price"`
	Stok      string `gorm:"type:varchar(255)" json:"stok"`
	IsActive  bool   `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *MenuUkuran) Prepare() error {
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
