package entity

import (
	"time"

	"gorm.io/gorm"
)

type OrderType struct {
	Id        int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code      string         `gorm:"type:varchar(20);unique;not null" json:"code"`
	Name      string         `gorm:"type:varchar(50);unique;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
