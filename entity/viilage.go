package entity

type Village struct {
	ID         int      `gorm:"primaryKey;column:id" json:"id"`
	DistrictID int      `gorm:"not null;column:district_id" json:"district_id"`
	Name       string   `gorm:"type:varchar(100);not null" json:"name"`
	Code       string   `gorm:"type:varchar(10);not null" json:"code"`
	FullCode   string   `gorm:"type:varchar(10);not null" json:"full_code"`
	PosCode    string   `gorm:"type:varchar(10);not null" json:"pos_code"`
	District   District `gorm:"foreignKey:DistrictID;references:ID" json:"district"`
}
