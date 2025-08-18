package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionItemAttribute struct {
	Id                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TransactionItemId  uuid.UUID
	ProductAttributeId uuid.UUID
	ProductAttribute   ProductAttribute
	AdditionalPrice    float64        `json:"additional_price"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
