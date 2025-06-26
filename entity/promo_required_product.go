package entity

type PromoRequiredProduct struct {
	PromoId   int `gorm:"primaryKey"`
	ProductId int `gorm:"primaryKey"`
}
