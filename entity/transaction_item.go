package entity

import (
	"time"
)

type TransactionItem struct {
	Id                 int                        `gorm:"primaryKey" json:"id"`
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
	UnitPrice          float64                    `json:"unit_price"`
	Price              float64                    `json:"price"`
	Discount           *float64                   `json:"discount"`
	Promo              *float64                   `json:"promo"`
	Rating             *float64                   `json:"rating"`
	Attributes         []TransactionItemAttribute `gorm:"foreignKey:TransactionItemId"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
