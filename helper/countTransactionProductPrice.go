package helper

import (
	"errors"
	"log"
	"time"

	"github.com/odhiahmad/kasirku-service/entity"
)

type TransactionPricing struct {
	TotalPrice float64
	BasePrice  float64
	Discount   float64
	Tax        float64
}

// HitungHargaTransaksi menentukan harga, diskon, dan promo yang berlaku saat transaksi
func HitungHargaTransaksi(
	product entity.Product,
	productVariantId *int,
	quantity int,
	allProductIds []int,
) (*TransactionPricing, error) {
	log.Printf("[HargaTransaksi] Menghitung harga untuk produk ID: %d, Qty: %d", product.Id, quantity)

	var price float64

	if product.HasVariant {
		if productVariantId == nil {
			return nil, errors.New("product variant ID is required for product with variants")
		}
		log.Printf("[HargaTransaksi] Produk memiliki variant. Mencari variant ID: %d", *productVariantId)

		var found bool
		for _, variant := range product.Variants {
			if variant.Id == *productVariantId {
				price = *variant.BasePrice
				log.Printf("[HargaTransaksi] Variant ditemukan. Base price: %.2f", price)
				found = true
				break
			}
		}
		if !found {
			log.Printf("[HargaTransaksi] Variant ID %d tidak ditemukan pada produk ini", *productVariantId)
			return nil, errors.New("variant not found for this product")
		}
	} else {
		price = *product.BasePrice
		log.Printf("[HargaTransaksi] Produk tanpa variant. Base price: %.2f", price)
	}

	var discount float64
	if product.Discount != nil && product.Discount.IsActive {
		now := time.Now()
		log.Printf("[HargaTransaksi] Diskon tersedia. ID: %d, Amount: %.2f, IsPercentage: %t, IsMultiple: %t", product.Discount.Id, product.Discount.Amount, product.Discount.IsPercentage, product.Discount.IsMultiple)
		log.Printf("[HargaTransaksi] StartAt: %s, EndAt: %s, Now: %s", product.Discount.StartAt.Format(time.RFC3339), product.Discount.EndAt.Format(time.RFC3339), now.Format(time.RFC3339))

		if (product.Discount.StartAt.IsZero() || now.After(product.Discount.StartAt)) &&
			(product.Discount.EndAt.IsZero() || now.Before(product.Discount.EndAt)) {

			var singleDiscount float64
			if product.Discount.IsPercentage {
				singleDiscount = price * product.Discount.Amount
			} else {
				singleDiscount = product.Discount.Amount
			}

			if product.Discount.IsMultiple {
				discount = singleDiscount * float64(quantity)
			} else {
				discount = singleDiscount
			}
			log.Printf("[HargaTransaksi] Diskon dihitung: %.2f", discount)
		} else {
			log.Printf("[HargaTransaksi] Diskon tidak berlaku untuk waktu sekarang.")
		}
	} else {
		log.Printf("[HargaTransaksi] Tidak ada diskon aktif untuk produk.")
	}

	var totalTax float64
	if product.Tax != nil && product.Tax.IsActive {
		if product.Tax.IsPercentage {
			totalTax += price * product.Tax.Amount
		} else {
			totalTax += product.Tax.Amount
		}
		log.Printf("[HargaTransaksi] Pajak dihitung: %.2f", totalTax)
	} else {
		log.Printf("[HargaTransaksi] Tidak ada pajak aktif untuk produk.")
	}

	totalPrice := price * float64(quantity)
	log.Printf("[HargaTransaksi] Total harga: %.2f, Total diskon: %.2f, Total pajak: %.2f", totalPrice, discount*float64(quantity), totalTax*float64(quantity))

	return &TransactionPricing{
		TotalPrice: totalPrice,
		BasePrice:  price,
		Discount:   discount,
		Tax:        totalTax * float64(quantity),
	}, nil
}
