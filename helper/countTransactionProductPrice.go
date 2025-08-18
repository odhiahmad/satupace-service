package helper

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/entity"
)

type TransactionPricing struct {
	TotalPrice float64
	SellPrice  float64
	BasePrice  float64
	Discount   float64
	Tax        float64
}

func HitungHargaTransaksi(
	product entity.Product,
	productVariantId *uuid.UUID,
	quantity int,
	allProductIds []uuid.UUID,
) (*TransactionPricing, error) {
	var price float64
	var sellPrice float64
	var basePrice float64

	if product.HasVariant {
		if productVariantId == nil {
			return nil, errors.New("product variant ID is required for product with variants")
		}

		var found bool
		for _, variant := range product.Variants {
			if variant.Id == *productVariantId {
				price = *variant.SellPrice
				sellPrice = *variant.SellPrice
				basePrice = *variant.BasePrice
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("variant not found for this product")
		}
	} else {
		price = *product.SellPrice
		sellPrice = *product.SellPrice
		basePrice = *product.BasePrice
	}

	var discount float64
	if product.Discount != nil && *product.Discount.IsActive {
		now := time.Now()

		if (product.Discount.StartAt.IsZero() || now.After(product.Discount.StartAt)) &&
			(product.Discount.EndAt.IsZero() || now.Before(product.Discount.EndAt)) {

			var singleDiscount float64
			if *product.Discount.IsPercentage {
				singleDiscount = price * product.Discount.Amount
			} else {
				singleDiscount = product.Discount.Amount
			}

			if *product.Discount.IsMultiple {
				discount = singleDiscount * float64(quantity)
			} else {
				discount = singleDiscount
			}
		}
	}

	var totalTax float64
	if product.Tax != nil {
		if *product.Tax.IsPercentage {
			totalTax += price * product.Tax.Amount
		} else {
			totalTax += product.Tax.Amount
		}
	}

	totalPrice := price * float64(quantity)

	return &TransactionPricing{
		TotalPrice: totalPrice,
		SellPrice:  sellPrice,
		BasePrice:  basePrice,
		Discount:   discount,
		Tax:        totalTax * float64(quantity),
	}, nil
}
