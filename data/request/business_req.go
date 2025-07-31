package request

import "github.com/google/uuid"

type BusinessCreate struct {
	Name           string  `json:"business_name" validate:"required"`
	OwnerName      string  `json:"owner_name" validate:"required"`
	BusinessTypeId int     `json:"business_type_id" validate:"required"`
	Image          *string `json:"image"`
	IsActive       bool    `json:"is_active" validate:"required"`
}

type BusinessUpdate struct {
	Id             uuid.UUID `json:"id" validate:"required"`
	Name           string    `json:"business_name" validate:"required"`
	OwnerName      string    `json:"owner_name" validate:"required"`
	BusinessTypeId int       `json:"business_type_id" validate:"required"`
	ProvinceID     int       `json:"province_id" validate:"required"`
	CityID         int       `json:"city_id" validate:"required"`
	DistrictID     int       `json:"district_id" validate:"required"`
	VillageID      int       `json:"village_id" validate:"required"`
	Image          *string   `json:"image"`
	IsActive       bool      `json:"is_active" validate:"required"`
}
