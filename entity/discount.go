package entity

import "time"

type Discount struct {
	Id           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessId   int       `gorm:"not null" json:"business_id"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	IsPercentage *bool     `gorm:"not null;default:false" json:"is_percentage"` // true = amount sebagai persen
	Amount       float64   `gorm:"not null" json:"amount"`                      // nilai diskon
	IsGlobal     *bool     `gorm:"not null;default:false" json:"is_global"`     // true = berlaku untuk semua produk
	IsMultiple   *bool     `gorm:"not null;default:false" json:"is_multiple"`   // true = berlaku kelipatan (misal beli 2x, diskon 2x)
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	IsActive     *bool     `gorm:"not null;default:false" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
