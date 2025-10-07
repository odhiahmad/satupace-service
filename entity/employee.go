package entity

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId  uuid.UUID
	Business    *Business `gorm:"foreignKey:BusinessId"`
	PhoneNumber *string   `gorm:"size:50" json:"phone_number"`
	Name        string    `gorm:"size:255" json:"name"`
	RoleId      int
	Role        *Role     `gorm:"foreignKey:RoleId"`
	Pin         string    `gorm:"not null" json:"-"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
