package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         *string   `gorm:"type:varchar(255)" json:"name"`
	Email        *string   `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	PendingEmail *string   `gorm:"type:varchar(255)" json:"pending_email"`
	PhoneNumber  string    `gorm:"type:varchar(255);uniqueIndex" json:"phone_number"`
	Gender       *string   `gorm:"type:varchar(255)" json:"gender"`
	Password     string    `gorm:"->;<-;not null" json:"-"`
	PinCode      string    `gorm:"type:varchar(255)" json:"-"`
	Token        string    `gorm:"-" json:"token"`
	HasProfile   bool      `gorm:"default:false" json:"has_profile"`
	IsVerified   bool      `gorm:"not null;column:is_verified" json:"is_verified"`
	IsActive     bool      `gorm:"default:false" json:"is_active"`
	IsSuspended  bool      `gorm:"default:false" json:"is_suspended"`
	ReportCount  int       `gorm:"default:0" json:"report_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
