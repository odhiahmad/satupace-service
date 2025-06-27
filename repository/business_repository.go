package repository

import (
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BusinessRepository interface {
	Create(business entity.Business) (entity.Business, error)
	Update(business entity.Business) (entity.Business, error)
	Delete(business entity.Business) error
	FindById(id int) (entity.Business, error)
	FindWithPagination(pagination request.Pagination) ([]entity.Business, int64, error)
}

type businessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) BusinessRepository {
	return &businessRepository{db}
}

// Create inserts a new business entity into the database.
func (r *businessRepository) Create(business entity.Business) (entity.Business, error) {
	err := r.db.Create(&business).Error
	return business, err
}

// Update modifies an existing business entity.
func (r *businessRepository) Update(business entity.Business) (entity.Business, error) {
	err := r.db.Save(&business).Error // Gunakan Save agar seluruh field diperbarui
	return business, err
}

// Delete removes a business entity.
func (r *businessRepository) Delete(business entity.Business) error {
	return r.db.Delete(&business).Error
}

// FindById retrieves a business entity by its ID, with branches and business type preloaded.
func (r *businessRepository) FindById(id int) (entity.Business, error) {
	var business entity.Business
	err := r.db.Preload("Branches").Preload("BusinessType").First(&business, id).Error
	return business, err
}

// FindWithPagination retrieves paginated business data, with optional search.
func (r *businessRepository) FindWithPagination(pagination request.Pagination) ([]entity.Business, int64, error) {
	var businesses []entity.Business
	var total int64

	// Base query
	baseQuery := r.db.Model(&entity.Business{}).
		Preload("Branches").
		Preload("BusinessType")

	// Apply search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR owner_name ILIKE ?", search, search)
	}

	// Count total data (tanpa pagination)
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Query utama dengan paginator
	_, _, err := p.Paginate(baseQuery, &businesses)
	if err != nil {
		return nil, 0, err
	}

	return businesses, total, nil
}
