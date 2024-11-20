package entity

import (
	"time"
)

type TransactionDetail struct {
	Id                 int `gorm:"type:int;primary_key"`
	TransactionID      uint
	Transaction        Transaction `gorm:"foreignKey:TransactionID"`
	ProductID          uint
	Product            Product `gorm:"foreignKey:ProductID"`
	ProductAttributeID uint
	ProductAttribute   ProductAttribute `gorm:"foreignKey:ProductAttributeID"`
	ProductSizeID      uint
	ProductSize        ProductSize `gorm:"foreignKey:ProductSizeID"`
	ProductVariantID   uint
	ProductVariant     ProductVariant `gorm:"foreignKey:ProductVariantID"`
	Total              string         `gorm:"type:varchar(255)" json:"total"`
	Discount           string         `gorm:"type:varchar(255)" json:"discount"`
	Promo              string         `gorm:"type:varchar(255)" json:"promo"`
	Rating             string         `gorm:"type:varchar(255)" json:"rating"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
