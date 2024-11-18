package entity

import (
	"time"
)

type User struct {
	Id         int      `gorm:"type:int;primary_key"`
	RoleId     int      `gorm:"column:role_id;not null" json:"role_id"`
	BusinessId int      `gorm:"column:business_id" json:"business_id"`
	Business   Business `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Email      string   `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password   string   `gorm:"->;<-;not null" json:"-"`
	Token      string   `gorm:"-" json:"token,omitempty"`
	IsActive   bool     `gorm:"not null; column:is_active"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
