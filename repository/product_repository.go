package repository

import (
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
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error // ✅ cukup satu pointer
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
	var bundles []entity.Product
	var total int64

	// Base query
	baseQuery := r.db.Model(&entity.Product{}).Where("business_id = ?", businessId)

	// Search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ? OR brand ILIKE ?", search, search)
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
