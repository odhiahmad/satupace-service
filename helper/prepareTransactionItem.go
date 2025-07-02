// helper/transaction.go
package helper

import (
	"errors"
	"time"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type TransactionItemInput struct {
	DB            *gorm.DB
	Items         []request.TransactionItemCreate // alias tipe dari item request
	AllProductIds []int
}

type TransactionItemInputUpdate struct {
	DB            *gorm.DB
	Items         []request.TransactionItemUpdate // alias tipe dari item request
	AllProductIds []int
}

type TransactionItemResult struct {
	Items         []entity.TransactionItem
	BasePrice     float64
	FinalPrice    float64
	TotalDiscount float64
	TotalTax      float64
	TotalPromo    float64
}

func PrepareTransactionItemsCreate(input TransactionItemInput) (TransactionItemResult, error) {
	var result TransactionItemResult

	for _, item := range input.Items {
		var attrs []entity.TransactionItemAttribute
		for _, attr := range item.Attributes {
			attrs = append(attrs, entity.TransactionItemAttribute{
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		productId := IntOrDefault(item.ProductId, 0)
		if productId != 0 {
			// Produk biasa
			var product entity.Product
			err := input.DB.Preload("Variants").
				Preload("Discount").
				Preload("Tax").
				First(&product, "id = ?", productId).Error
			if err != nil {
				return result, err
			}

			pricing, err := HitungHargaTransaksi(product, item.ProductVariantId, item.Quantity, input.AllProductIds)
			if err != nil {
				return result, err
			}

			subtotal := (pricing.TotalPrice + pricing.Tax) - pricing.Discount
			result.BasePrice += pricing.TotalPrice
			result.FinalPrice += subtotal
			result.TotalDiscount += pricing.Discount
			result.TotalTax += pricing.Tax

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          item.ProductId,
				BundleId:           nil,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				UnitPrice:          pricing.BasePrice,
				Total:              subtotal,
				Discount:           pricing.Discount,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
				Product:            &product, // untuk akses ulang promo minimum_spend
			})

		} else if item.BundleId != nil {
			// Bundle
			var bundle entity.Bundle
			err := input.DB.First(&bundle, "id = ?", item.BundleId).Error
			if err != nil {
				return result, err
			}

			pricing, err := HitungHargaBundle(bundle, item.Quantity)
			if err != nil {
				return result, err
			}

			subtotal := (pricing.Total + pricing.Tax) * float64(item.Quantity)
			result.BasePrice += pricing.BasePrice
			result.FinalPrice += subtotal
			result.TotalTax += pricing.Tax * float64(item.Quantity)

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          nil,
				BundleId:           item.BundleId,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				UnitPrice:          pricing.BasePrice,
				Total:              subtotal,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
			})
		} else {
			return result, errors.New("item harus memiliki product_id atau bundle_id")
		}
	}

	// ✅ Diskon global (setelah total dihitung)
	var globalDiscount entity.Discount
	err := input.DB.Where("is_global = ? AND is_active = ?", true, true).
		Where("start_at <= ? AND end_at >= ?", time.Now(), time.Now()).
		First(&globalDiscount).Error
	if err == nil {
		var totalDisc float64
		if globalDiscount.IsPercentage {
			totalDisc = result.FinalPrice * globalDiscount.Amount
		} else {
			totalDisc = globalDiscount.Amount
		}

		if globalDiscount.IsMultiple {
			var totalQty int
			for _, item := range result.Items {
				totalQty += item.Quantity
			}
			result.TotalDiscount += totalDisc * float64(totalQty)
		} else {
			result.TotalDiscount += totalDisc
		}
	}

	return result, nil
}

func PrepareTransactionItemsUpdate(input TransactionItemInputUpdate) (TransactionItemResult, error) {
	var result TransactionItemResult

	for _, item := range input.Items {
		var attrs []entity.TransactionItemAttribute
		for _, attr := range item.Attributes {
			attrs = append(attrs, entity.TransactionItemAttribute{
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		productId := IntOrDefault(item.ProductId, 0)
		if productId != 0 {
			// Produk biasa
			var product entity.Product
			err := input.DB.Preload("Variants").
				Preload("Discount").
				Preload("Tax").
				First(&product, "id = ?", productId).Error
			if err != nil {
				return result, err
			}

			pricing, err := HitungHargaTransaksi(product, item.ProductVariantId, item.Quantity, input.AllProductIds)
			if err != nil {
				return result, err
			}

			subtotal := (pricing.TotalPrice + pricing.Tax) - pricing.Discount
			result.BasePrice += pricing.TotalPrice
			result.FinalPrice += subtotal
			result.TotalDiscount += pricing.Discount
			result.TotalTax += pricing.Tax

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          item.ProductId,
				BundleId:           nil,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				UnitPrice:          pricing.BasePrice,
				Total:              subtotal,
				Discount:           pricing.Discount,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
				Product:            &product, // untuk akses ulang promo minimum_spend
			})

		} else if item.BundleId != nil {
			// Bundle
			var bundle entity.Bundle
			err := input.DB.First(&bundle, "id = ?", item.BundleId).Error
			if err != nil {
				return result, err
			}

			pricing, err := HitungHargaBundle(bundle, item.Quantity)
			if err != nil {
				return result, err
			}

			subtotal := (pricing.Total + pricing.Tax) * float64(item.Quantity)
			result.BasePrice += pricing.BasePrice
			result.FinalPrice += subtotal
			result.TotalTax += pricing.Tax * float64(item.Quantity)

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          nil,
				BundleId:           item.BundleId,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				UnitPrice:          pricing.BasePrice,
				Total:              subtotal,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
			})
		} else {
			return result, errors.New("item harus memiliki product_id atau bundle_id")
		}
	}

	// ✅ Diskon global (setelah total dihitung)
	var globalDiscount entity.Discount
	err := input.DB.Where("is_global = ? AND is_active = ?", true, true).
		Where("start_at <= ? AND end_at >= ?", time.Now(), time.Now()).
		First(&globalDiscount).Error
	if err == nil {
		var totalDisc float64
		if globalDiscount.IsPercentage {
			totalDisc = result.FinalPrice * globalDiscount.Amount
		} else {
			totalDisc = globalDiscount.Amount
		}

		if globalDiscount.IsMultiple {
			var totalQty int
			for _, item := range result.Items {
				totalQty += item.Quantity
			}
			result.TotalDiscount += totalDisc * float64(totalQty)
		} else {
			result.TotalDiscount += totalDisc
		}
	}

	return result, nil
}
