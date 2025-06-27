package request

type TaxCreate struct {
	BusinessId int     `json:"business_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	IsGlobal   bool    `json:"is_global"`
	ProductIds []int   `json:"product_ids,omitempty"` // hanya digunakan jika !IsGlobal
}

type TaxUpdate struct {
	Name       string  `json:"name" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	IsGlobal   bool    `json:"is_global"`
	ProductIds []int   `json:"product_ids,omitempty"` // untuk update relasi jika perlu
}
