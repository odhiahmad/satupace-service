package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductAttribute struct {
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProductId   uuid.UUID `gorm:"null" json:"product_id"`
	Product     *Product  `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Price       float64   `json:"price"`
	Image       *string   `gorm:"type:varchar(255)" json:"image"`
	IsAvailable *bool     `json:"is_available"`
	IsActive    *bool     `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
