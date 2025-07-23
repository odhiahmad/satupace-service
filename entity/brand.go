package entity

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	Id         int            `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessId int            `gorm:"not null" json:"business_id"`
	Business   *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
