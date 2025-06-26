package helper

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
)

type TransactionPricing struct {
	Price    float64
	Discount float64
	Promo    float64
	Tax      float64
}

// HitungHargaTransaksi menentukan harga, diskon, dan promo yang berlaku saat transaksi
func HitungHargaTransaksi(
	product entity.Product,
	productVariantId *int,
	quantity int,
	allProductIds []int,
	usedPromo map[int]bool, // ditambahkan parameter
) (*TransactionPricing, error) {
	var price float64

	// 1. Tentukan harga dari base atau variant
	if product.HasVariant {
		if productVariantId == nil {
			return nil, errors.New("product variant ID is required for product with variants")
		}
		var found bool
		for _, variant := range product.Variants {
			if variant.Id == *productVariantId {
				price = *variant.BasePrice
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("variant not found for this product")
		}
	} else {
		price = product.BasePrice
	}

	// 2. Hitung diskon
	var discount float64
	if product.Discount != nil && product.Discount.IsActive {
		if product.Discount.Type == "percent" {
			discount = price * (product.Discount.Amount / 100.0)
		} else {
			discount = product.Discount.Amount
		}
	}

	// 3. Hitung promo
	var promo float64
	for _, pp := range product.ProductPromos {
		p := pp.Promo
		if p == nil || !p.IsActive || usedPromo[p.Id] {
			continue
		}
		if quantity < pp.MinQuantity {
			continue
		}
		if len(p.RequiredProducts) > 0 && !containsAll(allProductIds, p.RequiredProducts) {
			continue
		}
		if p.Type == "percent" || p.Type == "percentage" {
			promo = price * (p.Amount / 100.0)
		} else {
			promo = p.Amount
		}
		usedPromo[p.Id] = true // tandai sudah dipakai
		break                  // hanya satu promo aktif
	}

	// 4. Hitung pajak
	var totalTax float64
	if product.Tax != nil && product.Tax.IsActive {
		if product.Tax.Type == "percentage" {
			totalTax += price * (product.Tax.Amount / 100.0)
		} else {
			totalTax += product.Tax.Amount
		}
	}

	return &TransactionPricing{
		Price:    price * float64(quantity),
		Discount: discount * float64(quantity),
		Promo:    promo * float64(quantity),
		Tax:      totalTax * float64(quantity),
	}, nil
}

func containsAll(cartProductIds []int, requiredProducts []entity.Product) bool {
	if len(requiredProducts) == 0 {
		return true
	}

	requiredMap := make(map[int]bool)
	for _, p := range requiredProducts {
		requiredMap[p.Id] = false
	}
	for _, id := range cartProductIds {
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
