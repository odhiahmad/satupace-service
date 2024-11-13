package entity

import (
	"time"
)

type User struct {
	Id        int `gorm:"type:int;primary_key"`
	RoleId    int
	Email     string `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password  string `gorm:"->;<-;not null" json:"-"`
	Token     string `gorm:"-" json:"token,omitempty"`
	IsActive  bool   `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
