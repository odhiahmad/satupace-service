package entity

type District struct {
	ID     int    `gorm:"primaryKey;column:id" json:"id"`
	CityID int    `gorm:"not null" json:"city_id"`
	Name   string `gorm:"type:varchar(255);not null" json:"name"`

	City City `gorm:"foreignKey:CityID;references:ID" json:"city"`
}
