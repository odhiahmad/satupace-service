package entity

type Village struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	DistrictID int    `gorm:"not null" json:"district_id"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`

	District District `gorm:"foreignKey:DistrictID;references:ID" json:"district"`
}
