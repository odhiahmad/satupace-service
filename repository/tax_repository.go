package repository

import (
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type TaxRepository interface {
	Create(tax entity.Tax) (entity.Tax, error)
	Update(tax entity.Tax) (entity.Tax, error)
	Delete(tax entity.Tax) error
	FindById(id int) (entity.Tax, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Tax, int64, error)
}

type taxRepository struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) TaxRepository {
	return &taxRepository{db}
}

func (r *taxRepository) Create(tax entity.Tax) (entity.Tax, error) {
	err := r.db.Create(&tax).Error
	return tax, err
}

func (r *taxRepository) Update(tax entity.Tax) (entity.Tax, error) {
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&tax).Error
	return tax, err
}

func (r *taxRepository) Delete(tax entity.Tax) error {
	return r.db.Delete(&tax).Error
}

func (r *taxRepository) FindById(id int) (entity.Tax, error) {
	var tax entity.Tax
	err := r.db.Preload("Products").First(&tax, id).Error
	return tax, err
}

func (r *taxRepository) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Tax, int64, error) {
	var tax []entity.Tax
	var total int64

	// Base query dengan preload relasi
	baseQuery := r.db.Model(&entity.Tax{}).
		Where("business_id = ?", businessId).
		Preload("Products")

	// Search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Hitung total sebelum paginasi
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &tax)
	if err != nil {
		return nil, 0, err
	}

	return tax, total, nil
}
