package entity

import (
	"time"

	"github.com/google/uuid"
)

type BundleItem struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BundleId  uuid.UUID `gorm:"index"`
	Bundle    Bundle    `gorm:"foreignKey:BundleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProductId uuid.UUID `gorm:"index"`
	Product   Product   `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity  int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
