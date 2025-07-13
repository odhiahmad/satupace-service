package entity

import (
	"time"
)

type UserBusiness struct {
	Id          int `gorm:"primaryKey;autoIncrement"`
	RoleId      int
	Role        Role `gorm:"foreignKey:RoleId"`
	BusinessId  int
	Business    Business    `gorm:"foreignKey:BusinessId"`
	Email       *string     `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	PhoneNumber string      `gorm:"type:varchar(255);uniqueIndex" json:"phone_number,omitempty"`
	Password    string      `gorm:"->;<-;not null" json:"-"`
	Token       string      `gorm:"-" json:"token,omitempty"`
	IsVerified  bool        `gorm:"not null; column:is_verified"`
	IsActive    bool        `gorm:"default:false" json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Membership  *Membership `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"membership,omitempty"`
}
