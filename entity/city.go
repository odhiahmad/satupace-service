package entity

type City struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	ProvinceID int    `gorm:"not null" json:"province_id"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`

	Province Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province"`
}
