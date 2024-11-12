package request

type MenuCreate struct {
	Nama        string `json:"nama" validate:"required"`
	Gambar      string `json:"gambar" validate:"required"`
	Price       string `json:"price" validate:"required"`
	Discount    string `json:"discount" validate:"required"`
	Promo       string `json:"promo" validate:"required"`
	Stok        string `json:"stok" validate:"required"`
	IsAvailable string `json:"is_available" validate:"required"`
}

type MenuUpdate struct {
	Id          int    `validate:"required"`
	Nama        string `json:"nama" validate:"required"`
	Gambar      string `json:"gambar" validate:"required"`
	Price       string `json:"price" validate:"required"`
	Discount    string `json:"discount" validate:"required"`
	Promo       string `json:"promo" validate:"required"`
	Stok        string `json:"stok" validate:"required"`
	IsAvailable string `json:"is_available" validate:"required"`
}
