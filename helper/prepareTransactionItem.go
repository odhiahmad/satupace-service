// helper/transaction.go
package helper

import (
	"errors"

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
	Total         float64
	TotalDiscount float64
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
			err := input.DB.Preload("Variants").Preload("Discount").Preload("ProductPromos.Promo").
				First(&product, "id = ?", productId).Error
			if err != nil {
				return result, err
			}

			pricing, err := HitungHargaTransaksi(product, item.ProductVariantId, item.Quantity, input.AllProductIds)
			if err != nil {
				return result, err
			}

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          item.ProductId,
				BundleId:           nil,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				Price:              pricing.Price,
				Discount:           pricing.Discount,
				Promo:              pricing.Promo,
				Rating:             item.Rating,
				Attributes:         attrs,
			})

			subtotal := (pricing.Price - pricing.Discount - pricing.Promo) * float64(item.Quantity)
			result.Total += subtotal
			result.TotalDiscount += pricing.Discount * float64(item.Quantity)
			result.TotalPromo += pricing.Promo * float64(item.Quantity)

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

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          nil,
				BundleId:           item.BundleId,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				Price:              pricing.Price,
				Discount:           pricing.Discount,
				Promo:              pricing.Promo,
				Rating:             item.Rating,
				Attributes:         attrs,
			})

			subtotal := (pricing.Price - pricing.Discount - pricing.Promo) * float64(item.Quantity)
			result.Total += subtotal
			result.TotalDiscount += pricing.Discount * float64(item.Quantity)
			result.TotalPromo += pricing.Promo * float64(item.Quantity)

		} else {
			return result, errors.New("item harus memiliki product_id atau bundle_id")
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
			err := input.DB.Preload("Variants").Preload("Discount").Preload("ProductPromos.Promo").
				First(&product, "id = ?", productId).Error
			if err != nil {
				return result, err
			}

			pricing, err := HitungHargaTransaksi(product, item.ProductVariantId, item.Quantity, input.AllProductIds)
			if err != nil {
				return result, err
			}

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          item.ProductId,
				BundleId:           nil,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				Price:              pricing.Price,
				Discount:           pricing.Discount,
				Promo:              pricing.Promo,
				Rating:             item.Rating,
				Attributes:         attrs,
			})

			subtotal := (pricing.Price - pricing.Discount - pricing.Promo) * float64(item.Quantity)
			result.Total += subtotal
			result.TotalDiscount += pricing.Discount * float64(item.Quantity)
			result.TotalPromo += pricing.Promo * float64(item.Quantity)

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

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          nil,
				BundleId:           item.BundleId,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				Price:              pricing.Price,
				Discount:           pricing.Discount,
				Promo:              pricing.Promo,
				Rating:             item.Rating,
				Attributes:         attrs,
			})

			subtotal := (pricing.Price - pricing.Discount - pricing.Promo) * float64(item.Quantity)
			result.Total += subtotal
			result.TotalDiscount += pricing.Discount * float64(item.Quantity)
			result.TotalPromo += pricing.Promo * float64(item.Quantity)

		} else {
			return result, errors.New("item harus memiliki product_id atau bundle_id")
		}
	}

	return result, nil
}
