package entity

import (
	"time"

	"gorm.io/gorm"
)

type BusinessBranch struct {
	gorm.Model
	BusinessID  uint
	Business    Business
	Pic         string `gorm:"type:varchar(255)" json:"pic"`
	PhoneNumber string `gorm:"type:varchar(255)" json:"phone_number"`
	Address     string `gorm:"type:varchar(255)" json:"address"`
	Lat         string `gorm:"type:varchar(255) null" json:"lat"`
	Long        string `gorm:"type:varchar(255) null" json:"long"`
	Rating      string `gorm:"type:varchar(255) null" json:"rating"`
	IsActive    bool   `gorm:"not null; column:is_active" json:"is_active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
