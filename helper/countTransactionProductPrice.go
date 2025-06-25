package helper

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
)

type TransactionPricing struct {
	Price    float64
	Discount float64
	Promo    float64
}

// HitungHargaTransaksi menentukan harga, diskon, dan promo yang berlaku saat transaksi
func HitungHargaTransaksi(
	product entity.Product,
	productVariantId *int,
	quantity int,
	allProductIds []int,
) (*TransactionPricing, error) {
	var price float64

	// 1. Tentukan harga
	if product.HasVariant {
		if productVariantId == nil {
			return nil, errors.New("product variant ID is required for product with variants")
		}
		var found bool
		for _, variant := range product.Variants {
			if variant.Id == *productVariantId {

				price = *variant.FinalPrice

				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("variant not found for this product")
		}
	} else {

		price = *product.FinalPrice

	}

	// 2. Hitung diskon
	var discount float64
	if product.Discount != nil {
		if product.Discount.Type == "percent" {
			discount = price * (product.Discount.Amount / 100.0)
		} else {
			discount = product.Discount.Amount
		}
	}

	// 3. Hitung promo (gunakan satu saja yang aktif, logikanya bisa dikembangkan)
	var promo float64
	for _, pp := range product.ProductPromos {
		if !pp.Promo.IsActive || quantity < pp.MinQuantity {
			continue
		}

		// ✅ Cek apakah RequiredProductIds terpenuhi
		if len(pp.Promo.RequiredProductIds) > 0 && !containsAll(allProductIds, pp.Promo.RequiredProductIds) {
			continue
		}

		// ✅ Hitung promo
		if pp.Promo.Type == "percent" {
			promo = price * (pp.Promo.Amount / 100.0)
		} else {
			promo = pp.Promo.Amount
		}
		break
	}

	return &TransactionPricing{
		Price:    price,
		Discount: discount,
		Promo:    promo,
	}, nil
}

func containsAll(target []int, required []int) bool {
	requiredMap := make(map[int]bool)
	for _, id := range required {
		requiredMap[id] = false
	}
	for _, id := range target {
		if _, ok := requiredMap[id]; ok {
			requiredMap[id] = true
		}
	}
	for _, found := range requiredMap {
		if !found {
			return false
		}
	}
	return true
}
