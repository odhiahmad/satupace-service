package entity

import "time"

type Tax struct {
	Id           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessId   int       `gorm:"not null" json:"business_id"`
	Business     *Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	IsPercentage bool      `gorm:"not null;default:false" json:"is_percentage"`
	IsGlobal     bool      `gorm:"not null;default:false" json:"is_global"`
	Amount       float64   `gorm:"not null" json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
