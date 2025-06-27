package repository

import (
	"time"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type ProductUnitRepository interface {
	Create(productUnit entity.ProductUnit) (entity.ProductUnit, error)
	Update(productUnit entity.ProductUnit) (entity.ProductUnit, error)
	Delete(id int) error
	FindById(id int) (entity.ProductUnit, error)
	FindActiveGlobalProductUnit(businessId int, now time.Time) (*entity.ProductUnit, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.ProductUnit, int64, error)
}

type productUnitRepo struct {
	db *gorm.DB
}

func NewProductUnitRepository(db *gorm.DB) ProductUnitRepository {
	return &productUnitRepo{db: db}
}

func (r *productUnitRepo) Create(productUnit entity.ProductUnit) (entity.ProductUnit, error) {
	err := r.db.Create(&productUnit).Error
	return productUnit, err
}

func (r *productUnitRepo) Update(productUnit entity.ProductUnit) (entity.ProductUnit, error) {
	err := r.db.Save(&productUnit).Error
	return productUnit, err
}

func (r *productUnitRepo) Delete(id int) error {
	return r.db.Delete(&entity.ProductUnit{}, id).Error
}

func (r *productUnitRepo) FindById(id int) (entity.ProductUnit, error) {
	var productUnit entity.ProductUnit
	err := r.db.Preload("Products").First(&productUnit, id).Error
	return productUnit, err
}

func (r *productUnitRepo) FindByBusinessId(businessId int) ([]entity.ProductUnit, error) {
	var productUnits []entity.ProductUnit
	err := r.db.Where("business_id = ?", businessId).Preload("Products").Find(&productUnits).Error
	return productUnits, err
}

func (r *productUnitRepo) FindActiveGlobalProductUnit(businessId int, now time.Time) (*entity.ProductUnit, error) {
	var productUnit entity.ProductUnit
	err := r.db.
		Where("business_id = ? AND is_global = ? AND start_at <= ? AND end_at >= ?", businessId, true, now, now).
		Order("start_at DESC").
		First(&productUnit).Error

	if err != nil {
		return nil, err
	}
	return &productUnit, nil
}

func (r *productUnitRepo) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.ProductUnit, int64, error) {
	var productUnits []entity.ProductUnit
	var total int64

	// Base query dengan preload relasi
	baseQuery := r.db.Model(&entity.ProductUnit{}).
		Where("business_id = ?", businessId).
		Preload("Products")

	// Filter pencarian
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Hitung total data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &productUnits)
	if err != nil {
		return nil, 0, err
	}

	return productUnits, total, nil
}
