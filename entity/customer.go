package entity

import (
	"time"
)

type Customer struct {
	Id          int    `gorm:"type:int;primary_key"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	PhoneNumber string `gorm:"uniqueIndex;type:varchar(255)" json:"phone_number"`
	IsActive    bool   `gorm:"not null" json:"is_active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
