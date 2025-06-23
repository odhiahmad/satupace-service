package entity

import "time"

type TransactionItemAttribute struct {
	Id                 int `gorm:"type:int;primary_key"`
	TransactionItemId  int
	ProductAttributeId int
	ProductAttribute   ProductAttribute
	AdditionalPrice    float64   `json:"additional_price"` // salin dari topping saat transaksi
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
