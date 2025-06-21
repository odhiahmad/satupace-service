package entity

type BundleItem struct {
	Id        int `gorm:"type:int;primary_key"`
	BundleId  int `gorm:"index"`
	ProductId int `gorm:"index"`
	Quantity  int
	Product   Product `gorm:"foreignKey:ProductId"`
}
