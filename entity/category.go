package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	Id         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BusinessId uuid.UUID      `gorm:"not null;index" json:"business_id"`
	Business   *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name       string         `gorm:"type:varchar(255)" json:"name"`
	ParentId   *uuid.UUID     `gorm:"index" json:"parent_id"`
	Children   []Category     `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"children"`
	IsActive   bool           `gorm:"not null" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
