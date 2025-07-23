package entity

import (
	"time"

	"gorm.io/gorm"
)

type TransactionItem struct {
	Id                 int                        `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionId      int                        `json:"transaction_id"`
	ProductId          *int                       `json:"product_id"`
	Product            *Product                   `gorm:"foreignKey:ProductId"`
	BundleId           *int                       `json:"bundle_id"`
	Bundle             *Bundle                    `gorm:"foreignKey:BundleId"`
	ProductAttributeId *int                       `json:"product_attribute_id"`
	ProductAttribute   *ProductAttribute          `gorm:"foreignKey:ProductAttributeId"`
	ProductVariantId   *int                       `json:"product_variant_id"`
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
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
