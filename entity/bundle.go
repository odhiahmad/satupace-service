package entity

import "time"

type Bundle struct {
	Id          int `gorm:"primaryKey"`
	BusinessId  int
	Business    Business `gorm:"foreignKey:BusinessId"`
	Name        string   `gorm:"type:varchar(255)"`
	Description *string  `gorm:"type:varchar(255)"`
	Image       *string  `gorm:"type:text"`
	BasePrice   float64
	Stock       int
	Items       []BundleItem `gorm:"foreignKey:BundleId"`
	IsAvailable bool         `gorm:"not null;column:is_available"`
	IsActive    bool         `gorm:"not null;column:is_active"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
