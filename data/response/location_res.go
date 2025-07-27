package response

type ProvinceResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CityResponse struct {
	ID         int    `json:"id"`
	ProvinceID int    `json:"province_id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
}

type DistrictResponse struct {
	ID       int    `json:"id"`
	CityID   int    `json:"city_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	FullCode string `json:"full_code"`
}

type VillageResponse struct {
	ID         int    `json:"id"`
	DistrictID int    `json:"district_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
	PosCode    string `json:"pos_code"`
}
