package entity

import "time"

type Discount struct {
	Id          int       `gorm:"primaryKey" json:"id"`
	BusinessId  int       `gorm:"not null" json:"business_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Type        string    `gorm:"type:varchar(20);not null" json:"type"` // "percentage", "fixed"
	Amount      float64   `gorm:"not null" json:"amount"`                // bisa berupa persen atau nominal tergantung kebutuhan
	IsGlobal    bool      `gorm:"not null" json:"is_global"`             // true = untuk semua produk
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	IsActive    bool      `gorm:"not null" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
