package repository

import (
	"log"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	Update(product entity.Product) (entity.Product, error)
	Delete(id int) error
	FindById(id int) (entity.Product, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error)
	SetActive(id int, active bool) error
	SetAvailable(id int, available bool) error
	SetHasVariant(productId int) error
	ResetVariantStateToFalse(productId int) error // ⬅️ Tambahkan ini
	WithTransaction(fn func(r ProductRepository) error) error
}

type productRepository struct {
	db *gorm.DB
}

func (r *productRepository) DB() *gorm.DB {
	return r.db
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error // ✅ cukup satu pointer
}

func (r *productRepository) Update(product entity.Product) (entity.Product, error) {
	// Kosongkan relasi agar GORM tidak abaikan update FK
	product.Tax = nil
	product.ProductPromos = nil
	product.Discount = nil
	product.ProductCategory = nil
	product.Unit = nil

	updateData := map[string]interface{}{
		"name":                product.Name,
		"description":         product.Description,
		"base_price":          product.BasePrice,
		"sku":                 product.SKU,
		"stock":               product.Stock,
		"minimum_sales":       *product.MinimumSales,
		"track_stock":         product.TrackStock,
		"has_variant":         product.HasVariant,
		"is_available":        product.IsAvailable,
		"is_active":           product.IsActive,
		"product_category_id": product.ProductCategoryId,
		"tax_id":              product.TaxId,
		"unit_id":             product.UnitId,
		"discount_id":         product.DiscountId,
	}

	log.Printf("Updating product with ID %d", product.Id)

	err := r.db.Debug().Model(&product).Updates(updateData).Error
	return product, err
}

func (r *productRepository) Delete(id int) error {
	var product entity.Product
	result := r.db.Where("id = ?", id).Delete(&product)
	return result.Error
}

func (r *productRepository) FindById(id int) (entity.Product, error) {
	var product entity.Product
	err := r.db.
		Preload("Variants").
		Preload("ProductCategory").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos").
		Preload("ProductPromos.Promo").
		Preload("ProductPromos.Promo.RequiredProducts").
		First(&product, id).Error
	return product, err
}

func (r *productRepository) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error) {
	var bundles []entity.Product
	var total int64

	// Base query
	baseQuery := r.db.Model(&entity.Product{}).
		Preload("Variants").
		Preload("ProductCategory").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos").
		Preload("ProductPromos.Promo").
		Preload("ProductPromos.Promo.RequiredProducts").
		Where("business_id = ?", businessId)

	// Search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ? OR brand ILIKE ?", search, search, search)
	}

	// Hitung total data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Gunakan helper paginator dengan validasi sort
	p := helper.Paginate(pagination)

	// Ambil data hasil paginasi
	_, _, err := p.Paginate(baseQuery, &bundles)
	if err != nil {
		return nil, 0, err
	}

	return bundles, total, nil
}

func (r *productRepository) SetActive(id int, active bool) error {
	return r.db.Model(&entity.Product{}).
		Where("id = ?", id).
		Update("is_active", active).Error
}

func (r *productRepository) SetAvailable(id int, available bool) error {
	return r.db.Model(&entity.Product{}).
		Where("id = ?", id).
		Update("is_available", available).Error
}

func (r *productRepository) SetHasVariant(productId int) error {

	updateData := map[string]interface{}{
		"has_variant": true,
		"base_price":  nil,
		"stock":       nil,
		"track_stock": false,
		"discount_id": nil,
		"promo_id":    nil,
	}

	return r.db.Model(&entity.Product{}).
		Where("id = ?", productId).
		Updates(updateData).Error

}

func (r *productRepository) ResetVariantStateToFalse(productId int) error {
	updateData := map[string]interface{}{
		"has_variant": false,
		"track_stock": true,
	}

	return r.db.Model(&entity.Product{}).
		Where("id = ?", productId).
		Updates(updateData).Error
}

func (r *productRepository) WithTransaction(fn func(r ProductRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &productRepository{db: tx}
		return fn(txRepo)
	})
}
