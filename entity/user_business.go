package entity

import (
	"time"
)

type UserBusiness struct {
	Id         int `gorm:"type:int;primary_key"`
	RoleId     int
	Role       Role `gorm:"foreignKey:RoleId"`
	BusinessId int
	Business   Business `gorm:"foreignKey:BusinessId"`
	Email      string   `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password   string   `gorm:"->;<-;not null" json:"-"`
	Token      string   `gorm:"-" json:"token,omitempty"`
	IsVerified bool     `gorm:"not null; column:is_verified"`
	IsActive   bool     `gorm:"not null; column:is_active"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
