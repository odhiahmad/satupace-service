package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserDeviceToken menyimpan FCM token per device per user.
// Satu user bisa punya banyak device (HP, tablet, dll).
type UserDeviceToken struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;index"`
	FCMToken  string    `gorm:"type:text;not null;uniqueIndex"`
	Platform  string    `gorm:"type:varchar(20)"` // android, ios, web
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
