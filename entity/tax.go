package entity

import "time"

type Tax struct {
	Id         int       `gorm:"primaryKey" json:"id"`
	BusinessId int       `gorm:"not null" json:"business_id"`
	Business   *Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Type       string    `gorm:"type:varchar(20);not null" json:"type"` // "percentage", "fixed"
	Amount     float64   `gorm:"not null" json:"amount"`
	IsGlobal   bool      `gorm:"default:false" json:"is_global"` // true = untuk semua produk
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	Products   []Product `gorm:"many2many:tax_relations;" json:"products"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
