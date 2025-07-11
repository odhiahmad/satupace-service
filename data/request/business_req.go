package request

type BusinessCreate struct {
	Name           string  `json:"business_name" binding:"required"`
	OwnerName      string  `json:"owner_name" binding:"required"`
	BusinessTypeId int     `json:"business_type_id" binding:"required"`
	Image          *string `json:"image,omitempty"`
	IsActive       bool    `json:"is_active" binding:"required"`
}

type BusinessUpdate struct {
	Id             int     `json:"id" binding:"required"`
	Name           string  `json:"business_name" binding:"required"`
	OwnerName      string  `json:"owner_name" binding:"required"`
	BusinessTypeId int     `json:"business_type_id" binding:"required"`
	ProvinceID     int     `json:"province_id" binding:"required"`
	CityID         int     `json:"city_id" binding:"required"`
	DistrictID     int     `json:"district_id" binding:"required"`
	VillageID      int     `json:"village_id" binding:"required"`
	Image          *string `json:"image,omitempty"`
	IsActive       bool    `json:"is_active" binding:"required"`
}
