package repository

import (
	"fmt"
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
	SetActive(id int, active bool) error
	SetAvailable(id int, available bool) error
	SetHasVariant(productId int) error
	ResetVariantStateToFalse(productId int) error // ⬅️ Tambahkan ini
	WithTransaction(fn func(conn ProductRepository) error) error
	UpdateImage(productId int, imageURL string) error
	IsSKUExist(sku string, businessId int) (bool, error)
	IsSKUExistExcept(sku string, businessId int, exceptProductId int) (bool, error)
	GetTx() *gorm.DB
	FindById(id int) (entity.Product, error)
	FindByIds(businessId int, ids []int) ([]entity.Product, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]entity.Product, string, error)
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
	err := conn.db.Create(&product).Error
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
		First(&product, product.Id).Error
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (conn *productConnection) Update(product entity.Product) (entity.Product, error) {
	product.Tax = nil
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
		First(&product, product.Id).Error

	return product, err
}

func (r *productConnection) UpdateAll(product *entity.Product) (entity.Product, error) {
	err := r.db.Model(&entity.Product{}).Where("id = ?", product.Id).Updates(product).Error
	if err != nil {
		return entity.Product{}, err
	}

	var updated entity.Product
	if err := r.db.
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
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

func (conn *productConnection) UpdateImage(productId int, imageURL string) error {
	return conn.db.Model(&entity.Product{}).
		Where("id = ?", productId).
		Updates(map[string]interface{}{
			"image":    imageURL,
			"is_ready": true,
		}).Error
}

func (conn *productConnection) IsSKUExist(sku string, businessId int) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Product{}).
		Where("business_id = ? AND sku = ?", businessId, sku).
		Count(&count).Error
	return count > 0, err
}

func (conn *productConnection) IsSKUExistExcept(sku string, businessId int, exceptProductId int) (bool, error) {
	var count int64
	err := conn.db.
		Model(&entity.Product{}).
		Where("business_id = ? AND sku = ? AND id != ?", businessId, sku, exceptProductId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (conn *productConnection) GetTx() *gorm.DB {
	return conn.db
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
		First(&product, id).Error
	return product, err
}

func (conn *productConnection) FindByIds(businessId int, ids []int) ([]entity.Product, error) {
	var products []entity.Product
	err := conn.db.
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Where("business_id = ? AND id IN ?", businessId, ids).
		Find(&products).Error

	return products, err
}

func (conn *productConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := conn.db.Model(&entity.Product{}).
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Where("business_id = ?", businessId)

	if pagination.CategoryID != nil {
		query = query.Where("category_id = ?", *pagination.CategoryID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})
	_, _, err := p.Paginate(query, &products)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (conn *productConnection) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]entity.Product, string, error) {
	var products []entity.Product

	query := conn.db.Model(&entity.Product{}).
		Preload("Variants").
		Preload("Brand").
		Preload("Category").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Where("business_id = ?", businessId)

	if pagination.CategoryID != nil {
		query = query.Where("category_id = ?", *pagination.CategoryID)
	}

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("products.name ILIKE ?", search)
	}

	sortBy := pagination.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "ASC"
	if pagination.OrderBy == "desc" {
		order = "DESC"
	}

	if pagination.Cursor != "" {
		cursorID, err := helper.DecodeCursorID(pagination.Cursor)
		if err != nil {
			return nil, "", err
		}

		if order == "ASC" {
			query = query.Where("products.id > ?", cursorID)
		} else {
			query = query.Where("products.id < ?", cursorID)
		}
	}

	limit := pagination.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query = query.Order(fmt.Sprintf("products.%s %s", sortBy, order)).Limit(limit + 1)

	if err := query.Find(&products).Error; err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(products) > limit {
		last := products[limit-1]
		nextCursor = helper.EncodeCursorID(int64(last.Id))
		products = products[:limit]
	}

	return products, nextCursor, nil
}
