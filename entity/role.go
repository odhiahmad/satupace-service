package entity

type Role struct {
	Id   int    `gorm:"type:int;primary_key"`
	Name string `gorm:"type:varchar(255)" json:"name"`
}
