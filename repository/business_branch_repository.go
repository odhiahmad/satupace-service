package repository

import (
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BusinessBranchRepository interface {
	Create(businessBranch entity.BusinessBranch) (entity.BusinessBranch, error)
	Update(businessBranch entity.BusinessBranch) (entity.BusinessBranch, error)
	Delete(businessBranch entity.BusinessBranch) error
	FindById(id int) (entity.BusinessBranch, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.BusinessBranch, int64, error)
}

type businessBranchConnection struct {
	db *gorm.DB
}

func NewBusinessBranchRepository(db *gorm.DB) BusinessBranchRepository {
	return &businessBranchConnection{db}
}

func (conn *businessBranchConnection) Create(businessBranch entity.BusinessBranch) (entity.BusinessBranch, error) {
	err := conn.db.Create(&businessBranch).Error
	return businessBranch, err
}

func (conn *businessBranchConnection) Update(businessBranch entity.BusinessBranch) (entity.BusinessBranch, error) {
	err := conn.db.Updates(&businessBranch).Error
	return businessBranch, err
}

func (conn *businessBranchConnection) Delete(businessBranch entity.BusinessBranch) error {
	return conn.db.Delete(&businessBranch).Error
}

func (conn *businessBranchConnection) FindById(id int) (entity.BusinessBranch, error) {
	var branch entity.BusinessBranch
	err := conn.db.First(&branch, id).Error
	return branch, err
}

func (conn *businessBranchConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.BusinessBranch, int64, error) {
	var bundles []entity.BusinessBranch
	var total int64

	// Base query untuk count
	baseQuery := conn.db.Model(&entity.Bundle{}).Where("business_id = ?", businessId)

	// Apply search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Hitung total data (tanpa cursor pagination)
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	// Query utama dengan paginator
	_, _, err := p.Paginate(baseQuery, &bundles)
	if err != nil {
		return nil, 0, err
	}

	return bundles, total, nil
}
