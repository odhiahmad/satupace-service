package entity

type City struct {
	ID         int      `gorm:"primaryKey;column:id" json:"id"`
	ProvinceID int      `gorm:"not null;column:province_id" json:"province_id"`
	Type       string   `gorm:"type:varchar(50);not null" json:"type"`
	Name       string   `gorm:"type:varchar(100);not null" json:"name"`
	Code       string   `gorm:"type:varchar(10);not null" json:"code"`
	FullCode   string   `gorm:"type:varchar(10);not null" json:"full_code"`
	Province   Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province"`
}
