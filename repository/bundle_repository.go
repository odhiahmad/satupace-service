package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BundleRepository interface {
	InsertBundle(bundle *entity.Bundle) error
	UpdateBundle(bundle *entity.Bundle) error
	FindById(bundleId int) (entity.Bundle, error)
	Delete(bundleId int) error
	InsertItemsByBundleId(bundleId int, items []entity.BundleItem) error
	DeleteItemsByBundleId(bundleId int) error
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Bundle, int64, error)
	SetIsActive(id int, isActive bool) error
}

type BundleConnection struct {
	Db *gorm.DB
}

func NewBundleRepository(Db *gorm.DB) BundleRepository {
	return &BundleConnection{Db: Db}
}

func (r *BundleConnection) InsertBundle(bundle *entity.Bundle) error {
	result := r.Db.Create(bundle)
	return result.Error
}

func (r *BundleConnection) UpdateBundle(bundle *entity.Bundle) error {
	result := r.Db.Save(bundle)
	return result.Error
}

func (r *BundleConnection) FindById(bundleId int) (entity.Bundle, error) {
	var bundle entity.Bundle
	result := r.Db.Preload("Items.Product").First(&bundle, bundleId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return bundle, errors.New("bundle product not found")
	}
	return bundle, result.Error
}

func (r *BundleConnection) Delete(bundleId int) error {
	if err := r.DeleteItemsByBundleId(bundleId); err != nil {
		return err
	}
	result := r.Db.Delete(&entity.Bundle{}, bundleId)
	return result.Error
}

func (r *BundleConnection) InsertItemsByBundleId(bundleId int, items []entity.BundleItem) error {
	for i := range items {
		items[i].BundleId = bundleId
	}
	result := r.Db.Create(&items)
	return result.Error
}

func (r *BundleConnection) DeleteItemsByBundleId(bundleId int) error {
	result := r.Db.Where("bundle_id = ?", bundleId).Delete(&entity.BundleItem{})
	return result.Error
}

func (r *BundleConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Bundle, int64, error) {
	var bundles []entity.Bundle
	var total int64

	// Base query untuk count
	baseQuery := r.Db.Model(&entity.Bundle{}).Preload("Items.Product").Where("business_id = ?", businessId)

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

func (r *BundleConnection) SetIsActive(id int, isActive bool) error {
	return r.Db.Model(&entity.Bundle{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}
