package entity

import "time"

type ProductUnit struct {
	Id         int       `gorm:"primaryKey" json:"id"`
	BusinessId int       `gorm:"not null" json:"business_id"`
	Business   *Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name       string    `gorm:"type:varchar(50);not null" json:"name"`
	Alias      string    `gorm:"type:varchar(20)" json:"alias"`
	Multiplier float64   `gorm:"default:1" json:"multiplier"`
	Products   []Product `gorm:"foreignKey:UnitId" json:"products"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
