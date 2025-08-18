package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserBusiness struct {
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         *string   `gorm:"type:varchar(255)" json:"name"`
	RoleId       int
	Role         Role `gorm:"foreignKey:RoleId"`
	BusinessId   uuid.UUID
	Business     Business    `gorm:"foreignKey:BusinessId"`
	Email        *string     `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	PendingEmail *string     `gorm:"type:varchar(255)" json:"pending_email"`
	PhoneNumber  string      `gorm:"type:varchar(255);uniqueIndex" json:"phone_number"`
	Password     string      `gorm:"->;<-;not null" json:"-"`
	PinCode      string      `gorm:"type:varchar(255)" json:"-"`
	Token        string      `gorm:"-" json:"token"`
	IsVerified   bool        `gorm:"not null; column:is_verified"`
	IsActive     bool        `gorm:"default:false" json:"is_active"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Membership   *Membership `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"membership"`
}
