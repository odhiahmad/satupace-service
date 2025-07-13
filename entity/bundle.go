package entity

import "time"

type Bundle struct {
	Id          int `gorm:"primaryKey"`
	BusinessId  int
	Business    Business `gorm:"foreignKey:BusinessId"`
	Name        string   `gorm:"type:varchar(255)"`
	Description *string  `gorm:"type:varchar(255)"`
	Image       *string  `gorm:"type:text"`
	BasePrice   *float64 `json:"base_price,omitempty"`
	SellPrice   *float64 `json:"sell_price,omitempty"`
	Stock       int
	Items       []BundleItem `gorm:"foreignKey:BundleId"`
	TaxId       *int         `gorm:"index" json:"tax_id,omitempty"`
	Tax         *Tax         `gorm:"foreignKey:TaxId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax,omitempty"`
	IsAvailable bool         `gorm:"not null;column:is_available"`
	IsActive    bool         `gorm:"not null;column:is_active"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
