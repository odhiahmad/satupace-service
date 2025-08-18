package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionItem struct {
	Id                 uuid.UUID                  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TransactionId      uuid.UUID                  `json:"transaction_id"`
	ProductId          *uuid.UUID                 `json:"product_id"`
	Product            *Product                   `gorm:"foreignKey:ProductId"`
	BundleId           *uuid.UUID                 `json:"bundle_id"`
	Bundle             *Bundle                    `gorm:"foreignKey:BundleId"`
	ProductAttributeId *uuid.UUID                 `json:"product_attribute_id"`
	ProductAttribute   *ProductAttribute          `gorm:"foreignKey:ProductAttributeId"`
	ProductVariantId   *uuid.UUID                 `json:"product_variant_id"`
	ProductVariant     *ProductVariant            `gorm:"foreignKey:ProductVariantId"`
	Quantity           int                        `json:"quantity"`
	BasePrice          float64                    `json:"basePrice"`
	SellPrice          float64                    `json:"sellPrice"`
	Total              float64                    `json:"total"`
	Discount           float64                    `json:"discount"`
	Promo              float64                    `json:"promo"`
	Tax                float64                    `json:"tax"`
	Rating             *float64                   `json:"rating"`
	Attributes         []TransactionItemAttribute `gorm:"foreignKey:TransactionItemId"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
