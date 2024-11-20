package entity

type ProductSubCategory struct {
	ID   int    `gorm:"type:int;primary_key"`
	Name string `gorm:"type:varchar(255)" json:"name"`
}
