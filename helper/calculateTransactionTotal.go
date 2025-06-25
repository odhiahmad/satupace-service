package helper

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type TransactionItemPricing struct {
	Price    float64 // Harga awal produk/bundle (setelah variant dan atribut jika ada)
	Discount float64 // Total diskon yang berlaku
	Promo    float64 // Total promo yang berlaku
	Subtotal float64 // (Price - Discount - Promo) * Quantity
}

type TransactionItemTotalResult struct {
	Items         []entity.TransactionItem
	Total         float64
	TotalDiscount float64
	TotalPromo    float64
}

func CalculateTransactionTotals(db *gorm.DB, items []entity.TransactionItem) (*TransactionItemTotalResult, error) {
	var total, totalDiscount, totalPromo float64

	for _, item := range items {
		var (
			pricing *TransactionPricing
			err     error
		)

		if item.ProductId != nil {
			var product entity.Product
			err = db.
				Preload("Variants").
				Preload("Discount").
				Preload("ProductPromos.Promo").
				First(&product, "id = ?", *item.ProductId).Error
			if err != nil {
				return nil, err
			}

			pricing, err = HitungHargaTransaksi(product, item.ProductVariantId, item.Quantity, nil)
			if err != nil {
				return nil, err
			}

		} else if item.BundleId != nil {
			var bundle entity.Bundle
			err = db.First(&bundle, "id = ?", *item.BundleId).Error
			if err != nil {
				return nil, err
			}

			bp, err := HitungHargaBundle(bundle, item.Quantity)
			if err != nil {
				return nil, err
			}
			pricing = ConvertBundleToTransactionPricing(bp)
		} else {
			continue // item tidak valid, skip
		}

		if pricing == nil {
			return nil, errors.New("gagal menghitung harga item transaksi")
		}

		subtotal := (pricing.Price - pricing.Discount - pricing.Promo) * float64(item.Quantity)
		total += subtotal
		totalDiscount += pricing.Discount * float64(item.Quantity)
		totalPromo += pricing.Promo * float64(item.Quantity)
	}

	return &TransactionItemTotalResult{
		Items:         items,
		Total:         total,
		TotalDiscount: totalDiscount,
		TotalPromo:    totalPromo,
	}, nil
}

func ConvertBundleToTransactionPricing(p *BundlePricing) *TransactionPricing {
	if p == nil {
		return nil
	}
	return &TransactionPricing{
		Price:    p.Price,
		Discount: p.Discount,
		Promo:    p.Promo,
	}
}
