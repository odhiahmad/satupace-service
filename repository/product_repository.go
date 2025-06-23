package repository

import (
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product entity.Product) (entity.Product, error)
	Update(product entity.Product) (entity.Product, error)
	Delete(id int) error
	FindById(id int) (entity.Product, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product entity.Product) (entity.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *productRepository) Update(product entity.Product) (entity.Product, error) {
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&product).Error
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
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos.Promo").
		First(&product, id).Error
	return product, err
}

func (r *productRepository) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	// Base query dengan preload relasi
	baseQuery := r.db.Model(&entity.Product{}).
		Where("business_id = ?", businessId).
		Preload("Variants").
		Preload("Tax").
		Preload("Discount").
		Preload("Unit").
		Preload("ProductPromos.Promo")

	// Filter pencarian
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where(
			"name ILIKE ? OR brand ILIKE ? OR description ILIKE ?",
			search, search, search,
		)
	}

	// Hitung total data sebelum paginasi
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Query dengan paginasi
	_, _, err := p.Paginate(baseQuery, &products)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
