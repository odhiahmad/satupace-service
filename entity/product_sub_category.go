package entity

type ProductSubCategory struct {
	ID                int `gorm:"type:int;primary_key"`
	ProductCategoryID uint
	ProductCategory   ProductCategory `gorm:"foreignKey:ProductCategoryID"`
	Name              string          `gorm:"type:varchar(255)" json:"name"`
}
