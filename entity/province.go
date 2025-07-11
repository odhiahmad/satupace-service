package entity

type Province struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	Code string `gorm:"type:varchar(10);not null" json:"code"`
}
