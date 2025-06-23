package response

type DiscountResponse struct {
	Id          int     `json:"id"`
	BusinessId  int     `json:"business_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	IsPercent   bool    `json:"is_percent"`
	IsGlobal    bool    `json:"is_global"`
	StartAt     string  `json:"start_at"`
	EndAt       string  `json:"end_at"`

	// Jika diskon spesifik, akan menampilkan daftar produk
	Products []ProductResponse `json:"products,omitempty"`
}
