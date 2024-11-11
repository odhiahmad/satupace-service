package request

type PerusahaanCreateDTO struct {
	Nama   string `json:"nama" validate:"required"`
	Alamat string `json:"alamat" validate:"required"`
	Lat    string `json:"lat" validate:"required"`
	Long   string `json:"long" validate:"required"`
	Logo   string `json:"logo" validate:"required"`
	Gambar string `json:"gambar" validate:"required"`
}

type PerusahaanUpdateDTO struct {
	Id     int    `validate:"required"`
	Nama   string `json:"nama" validate:"required"`
	Alamat string `json:"alamat" validate:"required"`
	Lat    string `json:"lat" validate:"required"`
	Long   string `json:"long" validate:"required"`
	Logo   string `json:"logo" validate:"required"`
	Gambar string `json:"gambar" validate:"required"`
}

type PerusahaanIsActiveDTO struct {
	Id       int
	IsActive bool `json:"is_active" binding:"required"`
}
