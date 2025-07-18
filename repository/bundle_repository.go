package repository

import (
	"errors"
	"fmt"

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
	SetIsActive(id int, isActive bool) error
	SetIsAvailable(id int, isActive bool) error
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Bundle, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]entity.Bundle, string, error)
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

func (conn *bundleConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Bundle, int64, error) {
	var bundles []entity.Bundle
	var total int64

	baseQuery := conn.db.Model(&entity.Bundle{}).Preload("Items.Product").Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &bundles)
	if err != nil {
		return nil, 0, err
	}

	return bundles, total, nil
}

func (conn *bundleConnection) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]entity.Bundle, string, error) {
	var bundles []entity.Bundle

	query := conn.db.Model(&entity.Bundle{}).
		Preload("Items.Product").
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	sortBy := pagination.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "ASC"
	if pagination.OrderBy == "desc" {
		order = "DESC"
	}

	if pagination.Cursor != "" {
		cursorID, err := helper.DecodeCursorID(pagination.Cursor)
		if err != nil {
			return nil, "", err
		}

		if order == "ASC" {
			query = query.Where("id > ?", cursorID)
		} else {
			query = query.Where("id < ?", cursorID)
		}
	}

	limit := pagination.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, order)).Limit(limit + 1)

	if err := query.Find(&bundles).Error; err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(bundles) > limit {
		last := bundles[limit-1]
		nextCursor = helper.EncodeCursorID(int64(last.Id))
		bundles = bundles[:limit]
	}

	return bundles, nextCursor, nil
}
