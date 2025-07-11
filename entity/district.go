package entity

type District struct {
	ID       int    `gorm:"primaryKey;column:id" json:"id"`
	CityID   int    `gorm:"not null;column:city_id" json:"city_id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Code     string `gorm:"type:varchar(10);not null" json:"code"`
	FullCode string `gorm:"type:varchar(10);not null" json:"full_code"`
	City     City   `gorm:"foreignKey:CityID;references:ID" json:"city"`
}
