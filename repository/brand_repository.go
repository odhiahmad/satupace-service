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

type BrandRepository interface {
	Create(brand entity.Brand) (entity.Brand, error)
	Update(brand entity.Brand) (entity.Brand, error)
	Delete(brand entity.Brand) error
	HasRelation(brandId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(brandId uuid.UUID) (brandes entity.Brand, err error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Brand, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Brand, string, bool, error)
}

type brandConnection struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandConnection{db}
}

func (conn *brandConnection) Create(brand entity.Brand) (entity.Brand, error) {
	err := conn.db.Create(&brand).Error
	if err != nil {
		return entity.Brand{}, err
	}

	err = conn.db.First(&brand, brand.Id).Error
	if err != nil {
		return entity.Brand{}, err
	}

	return brand, nil
}

func (conn *brandConnection) Update(brand entity.Brand) (entity.Brand, error) {
	err := conn.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&brand).Error
	if err != nil {
		return entity.Brand{}, err
	}

	err = conn.db.First(&brand, brand.Id).Error

	return brand, err
}

func (conn *brandConnection) Delete(brand entity.Brand) error {
	return conn.db.Delete(&brand).Error
}

func (conn *brandConnection) HasRelation(brandId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Product{}).Where("brand_id = ?", brandId).Count(&count).Error
	return count > 0, err
}

func (conn *brandConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Brand{}, id).Error
}

func (conn *brandConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.Brand{}, id).Error
}

func (conn *brandConnection) FindById(brandId uuid.UUID) (brandes entity.Brand, err error) {
	var brand entity.Brand
	result := conn.db.Find(&brand, brandId)
	if result != nil {
		return brand, nil
	} else {
		return brand, errors.New("tag is not found")
	}
}

func (conn *brandConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Brand, int64, error) {
	var brand []entity.Brand
	var total int64

	baseQuery := conn.db.Model(&entity.Brand{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &brand)
	if err != nil {
		return nil, 0, err
	}

	return brand, total, nil
}

func (conn *brandConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Brand, string, bool, error) {
	var brands []entity.Brand

	query := conn.db.Model(&entity.Brand{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ?", search)
	}

	sortBy := pagination.SortBy
	if sortBy == "" {
		sortBy = "updated_at"
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

	if err := query.Find(&brands).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(brands) > limit {
		last := brands[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		brands = brands[:limit]
		hasNext = true
	}

	return brands, nextCursor, hasNext, nil
}
