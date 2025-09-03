package repository

import (
	"errors"
	"fmt"

	"loka-kasir/data/request"
	"loka-kasir/entity"
	"loka-kasir/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BundleRepository interface {
	InsertBundle(bundle *entity.Bundle) (entity.Bundle, error)
	UpdateBundle(bundle *entity.Bundle) (entity.Bundle, error)
	FindById(bundleId uuid.UUID) (entity.Bundle, error)
	Delete(bundleId uuid.UUID) error
	InsertItemsByBundleId(bundleId uuid.UUID, items []entity.BundleItem) error
	DeleteItemsByBundleId(bundleId uuid.UUID) error
	SetIsActive(id uuid.UUID, isActive bool) error
	SetIsAvailable(id uuid.UUID, isActive bool) error
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Bundle, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Bundle, string, bool, error)
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

	err = conn.db.Preload("BundleItems").First(bundle, bundle.Id).Error
	return *bundle, err
}

func (conn *bundleConnection) UpdateBundle(bundle *entity.Bundle) (entity.Bundle, error) {
	err := conn.db.Save(bundle).Error
	if err != nil {
		return entity.Bundle{}, err
	}

	err = conn.db.Preload("BundleItems").First(bundle, bundle.Id).Error
	return *bundle, err
}

func (conn *bundleConnection) FindById(bundleId uuid.UUID) (entity.Bundle, error) {
	var bundle entity.Bundle
	result := conn.db.Preload("Items.Product").First(&bundle, bundleId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return bundle, errors.New("bundle product not found")
	}
	return bundle, result.Error
}

func (conn *bundleConnection) Delete(bundleId uuid.UUID) error {
	if err := conn.DeleteItemsByBundleId(bundleId); err != nil {
		return err
	}
	result := conn.db.Delete(&entity.Bundle{}, bundleId)
	return result.Error
}

func (conn *bundleConnection) InsertItemsByBundleId(bundleId uuid.UUID, items []entity.BundleItem) error {
	for i := range items {
		items[i].BundleId = bundleId
	}
	result := conn.db.Create(&items)
	return result.Error
}

func (conn *bundleConnection) DeleteItemsByBundleId(bundleId uuid.UUID) error {
	result := conn.db.Where("bundle_id = ?", bundleId).Delete(&entity.BundleItem{})
	return result.Error
}

func (conn *bundleConnection) SetIsActive(id uuid.UUID, isActive bool) error {
	return conn.db.Model(&entity.Bundle{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

func (conn *bundleConnection) SetIsAvailable(id uuid.UUID, isAvailable bool) error {
	return conn.db.Model(&entity.Bundle{}).
		Where("id = ?", id).
		Update("is_available", isAvailable).Error
}

func (conn *bundleConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Bundle, int64, error) {
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

func (conn *bundleConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Bundle, string, bool, error) {
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
			return nil, "", false, err
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
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(bundles) > limit {
		last := bundles[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		bundles = bundles[:limit]
		hasNext = true
	}

	return bundles, nextCursor, hasNext, nil
}
