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

type businessConnection struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) BusinessRepository {
	return &businessConnection{db}
}

// Create inserts a new business entity into the database.
func (conn *businessConnection) Create(business entity.Business) (entity.Business, error) {
	err := conn.db.Create(&business).Error
	return business, err
}

// Update modifies an existing business entity.
func (conn *businessConnection) Update(business entity.Business) (entity.Business, error) {
	err := conn.db.Save(&business).Error // Gunakan Save agar seluruh field diperbarui
	return business, err
}

// Delete removes a business entity.
func (conn *businessConnection) Delete(business entity.Business) error {
	return conn.db.Delete(&business).Error
}

// FindById retrieves a business entity by its ID, with branches and business type preloaded.
func (conn *businessConnection) FindById(id int) (entity.Business, error) {
	var business entity.Business
	err := conn.db.Preload("Branches").Preload("BusinessType").First(&business, id).Error
	return business, err
}

// FindWithPagination retrieves paginated business data, with optional search.
func (conn *businessConnection) FindWithPagination(pagination request.Pagination) ([]entity.Business, int64, error) {
	var businesses []entity.Business
	var total int64

	// Base query
	baseQuery := conn.db.Model(&entity.Business{}).
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
	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	// Query utama dengan paginator
	_, _, err := p.Paginate(baseQuery, &businesses)
	if err != nil {
		return nil, 0, err
	}

	return businesses, total, nil
}
