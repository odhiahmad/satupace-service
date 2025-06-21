package entity

import (
	"time"
)

type TransactionItem struct {
	Id                 int `gorm:"type:int;primary_key"`
	TransactionId      int
	ProductId          *int
	Product            Product `gorm:"foreignKey:ProductId"`
	BundleId           *int
	Bundle             Bundle `gorm:"foreignKey:BundleId"`
	ProductAttributeId *int
	ProductAttribute   ProductAttribute `gorm:"foreignKey:ProductAttributeId"`
	ProductVariantId   *int
	ProductVariant     ProductVariant `gorm:"foreignKey:ProductVariantId"`
	Quantity           int            `gorm:"type:varchar(255)" json:"quantity"`
	UnitPrice          float64        `gorm:"type:varchar(255)" json:"unit_price"`
	Price              float64        `gorm:"type:varchar(255)" json:"price"`
	Discount           float64        `gorm:"type:varchar(255)" json:"discount"`
	Promo              float64        `gorm:"type:varchar(255)" json:"promo"`
	Rating             float64        `gorm:"type:varchar(255)" json:"rating"`
	Attributes         []TransactionItemAttribute
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
