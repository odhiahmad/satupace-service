package repository

import (
	"log"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product entity.Product) (entity.Product, error)
	Update(product entity.Product) (entity.Product, error)
	UpdateAll(product *entity.Product) (entity.Product, error)
	Delete(id int) error
	FindById(id int) (entity.Product, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error)
	SetActive(id int, active bool) error
	SetAvailable(id int, available bool) error
	SetHasVariant(productId int) error
	ResetVariantStateToFalse(productId int) error // ⬅️ Tambahkan ini
	WithTransaction(fn func(conn ProductRepository) error) error
}

type productConnection struct {
	db *gorm.DB
}

func (conn *productConnection) DB() *gorm.DB {
	return conn.db
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productConnection{db}
}

func (conn *productConnection) Create(product entity.Product) (entity.Product, error) {
	// Buat data brand
	err := conn.db.Create(&product).Error
	if err != nil {
		return entity.Product{}, err
	}

	// Ambil ulang dengan preload Business
	err = conn.db.
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos").
		Preload("ProductPromos.Promo").
		Preload("ProductPromos.Promo.RequiredProducts").
		First(&product, product.Id).Error
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (conn *productConnection) Update(product entity.Product) (entity.Product, error) {
	// Kosongkan relasi agar GORM tidak abaikan update FK
	product.Tax = nil
	product.ProductPromos = nil
	product.Discount = nil
	product.Category = nil
	product.Unit = nil

	updateData := map[string]interface{}{
		"name":          product.Name,
		"description":   product.Description,
		"base_price":    product.BasePrice,
		"sell_price":    product.SellPrice,
		"sku":           product.SKU,
		"stock":         product.Stock,
		"minimum_sales": product.MinimumSales,
		"track_stock":   product.TrackStock,
		"has_variant":   product.HasVariant,
		"is_available":  product.IsAvailable,
		"is_active":     product.IsActive,
		"brand_id":      product.BrandId,
		"category_id":   product.CategoryId,
		"tax_id":        product.TaxId,
		"unit_id":       product.UnitId,
		"discount_id":   product.DiscountId,
	}

	log.Printf("Updating product with ID %d", product.Id)

	err := conn.db.Debug().Model(&product).Updates(updateData).Error
	if err != nil {
		return entity.Product{}, err
	}

	err = conn.db.
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos").
		Preload("ProductPromos.Promo").
		Preload("ProductPromos.Promo.RequiredProducts").
		First(&product, product.Id).Error

	return product, err
}

func (r *productConnection) UpdateAll(product *entity.Product) (entity.Product, error) {
	err := r.db.Model(&entity.Product{}).Where("id = ?", product.Id).Updates(product).Error
	if err != nil {
		return entity.Product{}, err
	}

	// Reload product dari DB setelah update
	var updated entity.Product
	if err := r.db.
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos").
		Preload("ProductPromos.Promo").
		Preload("ProductPromos.Promo.RequiredProducts").
		First(&updated, product.Id).Error; err != nil {
		return entity.Product{}, err
	}

	return updated, nil
}

func (conn *productConnection) Delete(id int) error {
	var product entity.Product
	result := conn.db.Where("id = ?", id).Delete(&product)
	return result.Error
}

func (conn *productConnection) FindById(id int) (entity.Product, error) {
	var product entity.Product
	err := conn.db.
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos").
		Preload("ProductPromos.Promo").
		Preload("ProductPromos.Promo.RequiredProducts").
		First(&product, id).Error
	return product, err
}

func (conn *productConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	// --- Gunakan Elasticsearch jika ada keyword pencarian ---
	if pagination.Search != "" {
		ids, err := helper.SearchProductElastic(pagination.Search, businessId)
		if err != nil {
			return nil, 0, err
		}

		// Fallback ke pencarian biasa jika hasil Elasticsearch kosong
		if len(ids) == 0 {
			baseQuery := conn.db.Model(&entity.Product{}).
				Preload("Variants").
				Preload("Brand").
				Preload("Category").
				Preload("Tax").
				Preload("Discount").
				Preload("Unit").
				Preload("ProductPromos").
				Preload("ProductPromos.Promo").
				Preload("ProductPromos.Promo.RequiredProducts").
				Where("business_id = ?", businessId).
				Where("name ILIKE ? OR description ILIKE ?", "%"+pagination.Search+"%", "%"+pagination.Search+"%")

			if err := baseQuery.Count(&total).Error; err != nil {
				return nil, 0, err
			}

			p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})
			_, _, err := p.Paginate(baseQuery, &products)
			if err != nil {
				return nil, 0, err
			}

			return products, total, nil
		}

		baseQuery := conn.db.Model(&entity.Product{}).
			Preload("Variants").
			Preload("Brand").
			Preload("Category").
			Preload("Tax").
			Preload("Discount").
			Preload("Unit").
			Preload("ProductPromos").
			Preload("ProductPromos.Promo").
			Preload("ProductPromos.Promo.RequiredProducts").
			Where("business_id = ?", businessId).
			Where("id IN ?", ids)

		if err := baseQuery.Count(&total).Error; err != nil {
			return nil, 0, err
		}

		p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})
		_, _, err = p.Paginate(baseQuery, &products)
		if err != nil {
			return nil, 0, err
		}

		return products, total, nil
	}

	// --- Default query tanpa search ---
	baseQuery := conn.db.Model(&entity.Product{}).
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos").
		Preload("ProductPromos.Promo").
		Preload("ProductPromos.Promo.RequiredProducts").
		Where("business_id = ?", businessId)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})
	_, _, err := p.Paginate(baseQuery, &products)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (conn *productConnection) SetActive(id int, active bool) error {
	return conn.db.Model(&entity.Product{}).
		Where("id = ?", id).
		Update("is_active", active).Error
}

func (conn *productConnection) SetAvailable(id int, available bool) error {
	return conn.db.Model(&entity.Product{}).
		Where("id = ?", id).
		Update("is_available", available).Error
}

func (conn *productConnection) SetHasVariant(productId int) error {

	updateData := map[string]interface{}{
		"has_variant": true,
		"base_price":  nil,
		"stock":       nil,
		"track_stock": false,
		"discount_id": nil,
		"promo_id":    nil,
	}

	return conn.db.Model(&entity.Product{}).
		Where("id = ?", productId).
		Updates(updateData).Error

}

func (conn *productConnection) ResetVariantStateToFalse(productId int) error {
	updateData := map[string]interface{}{
		"has_variant": false,
		"track_stock": true,
	}

	return conn.db.Model(&entity.Product{}).
		Where("id = ?", productId).
		Updates(updateData).Error
}

func (conn *productConnection) WithTransaction(fn func(conn ProductRepository) error) error {
	return conn.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &productConnection{db: tx}
		return fn(txRepo)
	})
}
