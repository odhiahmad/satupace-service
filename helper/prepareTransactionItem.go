// helper/transaction.go
package helper

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type TransactionItemInput struct {
	DB            *gorm.DB
	Items         []request.TransactionItemCreate // alias tipe dari item request
	AllProductIds []uuid.UUID
}

type TransactionItemInputUpdate struct {
	DB            *gorm.DB
	Items         []request.TransactionItemUpdate // alias tipe dari item request
	AllProductIds []uuid.UUID
}

type TransactionItemResult struct {
	Items         []entity.TransactionItem
	SellPrice     float64
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

		productId := item.ProductId
		if *productId != uuid.Nil {
			var product entity.Product
			err := input.DB.Preload("Variants").
				Preload("Discount").
				Preload("Tax").
				First(&product, "id = ?", productId).Error
			if err != nil {
				return result, err
			}

			if product.MinimumSales != nil {
				if item.Quantity < *product.MinimumSales {
					return result, errors.New("minimum pembelian untuk produk tidak terpenuhi")
				}
			}

			if item.ProductVariantId != nil && *item.ProductVariantId != uuid.Nil {
				for _, variant := range product.Variants {
					if variant.Id == *item.ProductVariantId {
						if variant.MinimumSales != nil && item.Quantity < *variant.MinimumSales {
							return result, errors.New("minimum pembelian untuk varian produk tidak terpenuhi")
						}
						break
					}
				}
			}

			pricing, err := HitungHargaTransaksi(product, item.ProductVariantId, item.Quantity, input.AllProductIds)
			if err != nil {
				return result, err
			}

			subtotal := (pricing.TotalPrice + pricing.Tax) - pricing.Discount
			result.SellPrice += (pricing.SellPrice * float64(item.Quantity))
			result.BasePrice += (pricing.BasePrice * float64(item.Quantity))
			result.FinalPrice += subtotal
			result.TotalDiscount += pricing.Discount
			result.TotalTax += pricing.Tax

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          item.ProductId,
				BundleId:           nil,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				BasePrice:          pricing.BasePrice,
				SellPrice:          pricing.SellPrice,
				Total:              subtotal,
				Discount:           pricing.Discount,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
				Product:            &product,
			})

		} else if item.BundleId != nil {
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
			result.SellPrice += (pricing.SellPrice * float64(item.Quantity))
			result.BasePrice += (pricing.BasePrice * float64(item.Quantity))
			result.FinalPrice += subtotal
			result.TotalTax += pricing.Tax * float64(item.Quantity)

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          nil,
				BundleId:           item.BundleId,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				BasePrice:          pricing.BasePrice,
				SellPrice:          pricing.SellPrice,
				Total:              subtotal,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
			})
		} else {
			return result, errors.New("item harus memiliki product_id atau bundle_id")
		}
	}

	var globalDiscount entity.Discount
	err := input.DB.Where("is_global = ? AND is_active = ?", true, true).
		Where("start_at <= ? AND end_at >= ?", time.Now(), time.Now()).
		First(&globalDiscount).Error
	if err == nil {
		var totalDisc float64
		if *globalDiscount.IsPercentage {
			totalDisc = result.FinalPrice * globalDiscount.Amount
		} else {
			totalDisc = globalDiscount.Amount
		}

		if *globalDiscount.IsMultiple {
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

		productId := item.ProductId
		if *productId != uuid.Nil {
			var product entity.Product
			err := input.DB.Preload("Variants").
				Preload("Discount").
				Preload("Tax").
				First(&product, "id = ?", productId).Error
			if err != nil {
				return result, err
			}

			if product.MinimumSales != nil {
				if item.Quantity < *product.MinimumSales {
					return result, errors.New("minimum pembelian untuk produk tidak terpenuhi")
				}
			}

			if item.ProductVariantId != nil && *item.ProductVariantId != uuid.Nil {
				for _, variant := range product.Variants {
					if variant.Id == *item.ProductVariantId {
						if variant.MinimumSales != nil && item.Quantity < *variant.MinimumSales {
							return result, errors.New("minimum pembelian untuk varian produk tidak terpenuhi")
						}
						break
					}
				}
			}

			pricing, err := HitungHargaTransaksi(product, item.ProductVariantId, item.Quantity, input.AllProductIds)
			if err != nil {
				return result, err
			}

			subtotal := (pricing.TotalPrice + pricing.Tax) - pricing.Discount
			result.SellPrice += (pricing.SellPrice * float64(item.Quantity))
			result.BasePrice += (pricing.BasePrice * float64(item.Quantity))
			result.FinalPrice += subtotal
			result.TotalDiscount += pricing.Discount
			result.TotalTax += pricing.Tax

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          item.ProductId,
				BundleId:           nil,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				BasePrice:          pricing.BasePrice,
				SellPrice:          pricing.SellPrice,
				Total:              subtotal,
				Discount:           pricing.Discount,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
				Product:            &product,
			})

		} else if item.BundleId != nil {
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
			result.SellPrice += (pricing.SellPrice * float64(item.Quantity))
			result.BasePrice += (pricing.BasePrice * float64(item.Quantity))
			result.FinalPrice += subtotal
			result.TotalTax += pricing.Tax * float64(item.Quantity)

			result.Items = append(result.Items, entity.TransactionItem{
				ProductId:          nil,
				BundleId:           item.BundleId,
				ProductAttributeId: item.ProductAttributeId,
				ProductVariantId:   item.ProductVariantId,
				Quantity:           item.Quantity,
				BasePrice:          pricing.BasePrice,
				SellPrice:          pricing.SellPrice,
				Total:              subtotal,
				Tax:                pricing.Tax,
				Rating:             item.Rating,
				Attributes:         attrs,
			})
		} else {
			return result, errors.New("item harus memiliki product_id atau bundle_id")
		}
	}

	var globalDiscount entity.Discount
	err := input.DB.Where("is_global = ? AND is_active = ?", true, true).
		Where("start_at <= ? AND end_at >= ?", time.Now(), time.Now()).
		First(&globalDiscount).Error
	if err == nil {
		var totalDisc float64
		if *globalDiscount.IsPercentage {
			totalDisc = result.FinalPrice * globalDiscount.Amount
		} else {
			totalDisc = globalDiscount.Amount
		}

		if *globalDiscount.IsMultiple {
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
