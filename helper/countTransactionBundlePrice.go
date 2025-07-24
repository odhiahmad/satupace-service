package helper

import (
	"log"

	"github.com/odhiahmad/kasirku-service/entity"
)

type BundlePricing struct {
	Total     float64
	SellPrice float64
	BasePrice float64
	Tax       float64
}

func HitungHargaBundle(bundle entity.Bundle, quantity int) (*BundlePricing, error) {
	// Hitung harga dasar berdasarkan harga satuan bundle dikali kuantitas
	price := *bundle.SellPrice * float64(quantity)
	var totalTax float64

	if bundle.Tax != nil {
		if *bundle.Tax.IsPercentage {
			totalTax = price * (bundle.Tax.Amount / 100.0)
		} else {
			totalTax = bundle.Tax.Amount * float64(quantity)
		}
		log.Printf("[HitungHargaBundle] Pajak dihitung: %.2f", totalTax)
	} else {
		log.Printf("[HitungHargaBundle] Tidak ada pajak aktif untuk bundle.")
	}

	total := price + totalTax

	return &BundlePricing{
		Total:     total,
		SellPrice: *bundle.SellPrice,
		BasePrice: *bundle.BasePrice,
		Tax:       totalTax,
	}, nil
}
