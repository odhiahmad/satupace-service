package helper

import (
	"github.com/odhiahmad/kasirku-service/entity"
)

type BundlePricing struct {
	Price    float64
	Discount float64
	Promo    float64
}

// HitungHargaBundle menghitung harga, diskon, dan promo dari bundle
func HitungHargaBundle(bundle entity.Bundle, quantity int) (*BundlePricing, error) {
	// Saat ini hanya menggunakan BasePrice, diskon dan promo belum diterapkan
	return &BundlePricing{
		Price:    bundle.BasePrice,
		Discount: 0,
		Promo:    0,
	}, nil
}
