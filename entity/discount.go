package entity

import "time"

type Discount struct {
	Id          int       `gorm:"primaryKey" json:"id"`
	BusinessId  int       `gorm:"not null" json:"business_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Type        string    `gorm:"type:discount_type;not null" json:"type"` // "percent" atau "fixed"
	Amount      float64   `gorm:"not null" json:"amount"`                  // bisa berupa persen atau nominal tergantung kebutuhan
	IsPercent   bool      `gorm:"not null" json:"is_percent"`
	IsGlobal    bool      `gorm:"not null" json:"is_global"` // true = untuk semua produk
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Products    []Product `gorm:"many2many:product_discounts;" json:"products,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
