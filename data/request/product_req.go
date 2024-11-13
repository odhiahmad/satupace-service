package request

type ProductCreate struct {
	Name        string `json:"nama" validate:"required"`
	Image       string `json:"gambar" validate:"required"`
	Price       uint   `json:"price" validate:"required"`
	IsAvailable bool   `json:"is_available" validate:"required"`
}

type ProductUpdate struct {
	Id          int    `validate:"required"`
	Nama        string `json:"nama" validate:"required"`
	Gambar      string `json:"gambar" validate:"required"`
	Price       string `json:"price" validate:"required"`
	Discount    string `json:"discount" validate:"required"`
	Promo       string `json:"promo" validate:"required"`
	Stok        string `json:"stok" validate:"required"`
	IsAvailable bool   `json:"is_available" validate:"required"`
}
