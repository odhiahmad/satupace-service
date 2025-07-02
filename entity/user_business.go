package entity

import (
	"time"
)

type UserBusiness struct {
	Id          int `gorm:"type:int;primary_key"`
	RoleId      int
	Role        Role `gorm:"foreignKey:RoleId"`
	BusinessId  int
	Business    Business        `gorm:"foreignKey:BusinessId"`
	BranchId    *int            `gorm:"index"`
	Branch      *BusinessBranch `gorm:"foreignKey:BranchId"`
	Email       string          `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	PhoneNumber *string         `gorm:"type:varchar(255)" json:"phone_number,omitempty"`
	Password    string          `gorm:"->;<-;not null" json:"-"`
	Token       string          `gorm:"-" json:"token,omitempty"`
	IsVerified  bool            `gorm:"not null; column:is_verified"`
	IsActive    bool            `gorm:"not null; column:is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Memberships []Membership    `gorm:"foreignKey:UserId" json:"memberships,omitempty"`
}
