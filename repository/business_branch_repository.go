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

type businessBranchRepository struct {
	db *gorm.DB
}

func NewBusinessBranchRepository(db *gorm.DB) BusinessBranchRepository {
	return &businessBranchRepository{db}
}

func (r *businessBranchRepository) Create(businessBranch entity.BusinessBranch) (entity.BusinessBranch, error) {
	err := r.db.Create(&businessBranch).Error
	return businessBranch, err
}

func (r *businessBranchRepository) Update(businessBranch entity.BusinessBranch) (entity.BusinessBranch, error) {
	err := r.db.Updates(&businessBranch).Error
	return businessBranch, err
}

func (r *businessBranchRepository) Delete(businessBranch entity.BusinessBranch) error {
	return r.db.Delete(&businessBranch).Error
}

func (r *businessBranchRepository) FindById(id int) (entity.BusinessBranch, error) {
	var branch entity.BusinessBranch
	err := r.db.First(&branch, id).Error
	return branch, err
}

func (r *businessBranchRepository) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.BusinessBranch, int64, error) {
	var bundles []entity.BusinessBranch
	var total int64

	// Base query untuk count
	baseQuery := r.db.Model(&entity.Bundle{}).Where("business_id = ?", businessId)

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
	p := helper.Paginate(pagination)

	// Query utama dengan paginator
	_, _, err := p.Paginate(baseQuery, &bundles)
	if err != nil {
		return nil, 0, err
	}

	return bundles, total, nil
}
