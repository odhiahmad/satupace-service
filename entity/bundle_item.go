package entity

import "time"

type BundleItem struct {
	Id        int `gorm:"primaryKey;autoIncrement" json:"id"`
	BundleId  int `gorm:"index"`
	ProductId int `gorm:"index"`
	Quantity  int
	Product   Product   `gorm:"foreignKey:ProductId"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
