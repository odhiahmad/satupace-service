package request

type ProvinceRequest struct {
	ProvinceID int `form:"province_id" json:"province_id" validate:"required"`
}

type CityRequest struct {
	CityID int `form:"city_id" json:"city_id" validate:"required"`
}

type DistrictRequest struct {
	DistrictID int `form:"district_id" json:"district_id" validate:"required"`
}
