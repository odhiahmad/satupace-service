package entity

import (
	"time"
)

type User struct {
	Id        int       `gorm:"type:int;primary_key"`
	RoleID    int       `gorm:"column:role_id;not null" json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID"`
	Email     string    `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password  string    `gorm:"->;<-;not null" json:"-"`
	Token     string    `gorm:"-" json:"token,omitempty"`
	IsActive  bool      `gorm:"not null; column:is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
