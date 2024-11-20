package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserBusiness struct {
	gorm.Model
	RoleID     int
	Role       Role `gorm:"foreignKey:RoleID"`
	BusinessID uint
	Business   Business `gorm:"foreignKey:BusinessID"`
	Email      string   `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password   string   `gorm:"->;<-;not null" json:"-"`
	Token      string   `gorm:"-" json:"token,omitempty"`
	IsVerified bool     `gorm:"not null; column:is_verified"`
	IsActive   bool     `gorm:"not null; column:is_active"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
