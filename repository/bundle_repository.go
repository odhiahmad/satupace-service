package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BundleRepository interface {
	InsertBundle(bundle *entity.Bundle) (entity.Bundle, error)
	UpdateBundle(bundle *entity.Bundle) (entity.Bundle, error)
	FindById(bundleId int) (entity.Bundle, error)
	Delete(bundleId int) error
	InsertItemsByBundleId(bundleId int, items []entity.BundleItem) error
	DeleteItemsByBundleId(bundleId int) error
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Bundle, int64, error)
	SetIsActive(id int, isActive bool) error
	SetIsAvailable(id int, isActive bool) error
}

type bundleConnection struct {
	db *gorm.DB
}

func NewBundleRepository(db *gorm.DB) BundleRepository {
	return &bundleConnection{db: db}
}

func (conn *bundleConnection) InsertBundle(bundle *entity.Bundle) (entity.Bundle, error) {
	err := conn.db.Create(bundle).Error
	if err != nil {
		return entity.Bundle{}, err
	}

	// Misal ingin preload relasi seperti BundleItems atau lainnya
	err = conn.db.Preload("BundleItems").First(bundle, bundle.Id).Error
	return *bundle, err
}

func (conn *bundleConnection) UpdateBundle(bundle *entity.Bundle) (entity.Bundle, error) {
	err := conn.db.Save(bundle).Error
	if err != nil {
		return entity.Bundle{}, err
	}

	// Preload relasi jika perlu
	err = conn.db.Preload("BundleItems").First(bundle, bundle.Id).Error
	return *bundle, err
}

func (conn *bundleConnection) FindById(bundleId int) (entity.Bundle, error) {
	var bundle entity.Bundle
	result := conn.db.Preload("Items.Product").First(&bundle, bundleId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return bundle, errors.New("bundle product not found")
	}
	return bundle, result.Error
}

func (conn *bundleConnection) Delete(bundleId int) error {
	if err := conn.DeleteItemsByBundleId(bundleId); err != nil {
		return err
	}
	result := conn.db.Delete(&entity.Bundle{}, bundleId)
	return result.Error
}

func (conn *bundleConnection) InsertItemsByBundleId(bundleId int, items []entity.BundleItem) error {
	for i := range items {
		items[i].BundleId = bundleId
	}
	result := conn.db.Create(&items)
	return result.Error
}

func (conn *bundleConnection) DeleteItemsByBundleId(bundleId int) error {
	result := conn.db.Where("bundle_id = ?", bundleId).Delete(&entity.BundleItem{})
	return result.Error
}

func (conn *bundleConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Bundle, int64, error) {
	var bundles []entity.Bundle
	var total int64

	// Base query untuk count
	baseQuery := conn.db.Model(&entity.Bundle{}).Preload("Items.Product").Where("business_id = ?", businessId)

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

func (conn *bundleConnection) SetIsActive(id int, isActive bool) error {
	return conn.db.Model(&entity.Bundle{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

func (conn *bundleConnection) SetIsAvailable(id int, isAvailable bool) error {
	return conn.db.Model(&entity.Bundle{}).
		Where("id = ?", id).
		Update("is_available", isAvailable).Error
}
