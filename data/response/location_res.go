package response

type ProvinceResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"` // ✅ ditambahkan
}

type CityResponse struct {
	ID         int    `json:"id"`
	ProvinceID int    `json:"province_id"`
	Type       string `json:"type"` // ✅ ditambahkan
	Name       string `json:"name"`
	Code       string `json:"code"`      // ✅ ditambahkan
	FullCode   string `json:"full_code"` // ✅ ditambahkan
}

type DistrictResponse struct {
	ID       int    `json:"id"`
	CityID   int    `json:"city_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`      // ✅ ditambahkan
	FullCode string `json:"full_code"` // ✅ ditambahkan
}

type VillageResponse struct {
	ID         int    `json:"id"`
	DistrictID int    `json:"district_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`      // ✅ ditambahkan
	FullCode   string `json:"full_code"` // ✅ ditambahkan
	PosCode    string `json:"pos_code"`  // ✅ ditambahkan
}
