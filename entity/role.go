package entity

type Role struct {
	Id   int    `gorm:"type:int;primary_key"`
	Nama string `gorm:"type:varchar(255)" json:"nama"`
}
