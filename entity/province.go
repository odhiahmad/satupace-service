package entity

type Province struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}
