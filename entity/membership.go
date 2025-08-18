package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Membership struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    uuid.UUID      `gorm:"not null" json:"user_id"`
	User      *UserBusiness  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StartDate time.Time      `gorm:"not null" json:"start_date"`
	EndDate   time.Time      `gorm:"not null" json:"end_date"`
	Type      string         `gorm:"type:varchar(20);not null" json:"type"`
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
