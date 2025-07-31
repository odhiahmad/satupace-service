package entity

import (
	"time"

	"github.com/google/uuid"
)

type Business struct {
	Id             uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name           string        `gorm:"type:varchar(255);not null" json:"business_name"`
	OwnerName      string        `gorm:"type:varchar(255);not null" json:"owner_name"`
	BusinessTypeId *int          `gorm:"not null" json:"business_type_id"`
	BusinessType   *BusinessType `gorm:"foreignKey:BusinessTypeId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProvinceID     *int          `json:"province_id"`
	CityID         *int          `json:"city_id"`
	DistrictID     *int          `json:"district_id"`
	VillageID      *int          `json:"village_id"`
	Province       *Province     `gorm:"foreignKey:ProvinceID;references:ID" json:"province"`
	City           *City         `gorm:"foreignKey:CityID;references:ID" json:"city"`
	District       *District     `gorm:"foreignKey:DistrictID;references:ID" json:"district"`
	Village        *Village      `gorm:"foreignKey:VillageID;references:ID" json:"village"`
	Image          *string       `gorm:"type:varchar(255)" json:"image"`
	IsActive       bool          `gorm:"not null" json:"is_active"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
