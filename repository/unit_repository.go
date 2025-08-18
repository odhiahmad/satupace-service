package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type UnitRepository interface {
	Create(unit entity.Unit) (entity.Unit, error)
	Update(unit entity.Unit) (entity.Unit, error)
	Delete(id uuid.UUID) error
	HasRelation(brandId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(id uuid.UUID) (entity.Unit, error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Unit, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Unit, string, bool, error)
}

type unitConnection struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitConnection{db: db}
}

func (connection *unitConnection) Create(unit entity.Unit) (entity.Unit, error) {
	err := connection.db.Create(&unit).Error
	return unit, err
}

func (connection *unitConnection) Update(unit entity.Unit) (entity.Unit, error) {
	err := connection.db.Save(&unit).Error
	return unit, err
}

func (connection *unitConnection) Delete(id uuid.UUID) error {
	return connection.db.Delete(&entity.Unit{}, id).Error
}

func (conn *unitConnection) HasRelation(unitId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Product{}).Where("unit_id = ?", unitId).Count(&count).Error
	return count > 0, err
}

func (conn *unitConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Unit{}, id).Error
}

func (conn *unitConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.Unit{}, id).Error
}

func (connection *unitConnection) FindById(id uuid.UUID) (entity.Unit, error) {
	var unit entity.Unit
	err := connection.db.First(&unit, id).Error
	return unit, err
}

func (connection *unitConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Unit, int64, error) {
	var units []entity.Unit
	var total int64

	baseQuery := connection.db.Model(&entity.Unit{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &units)
	if err != nil {
		return nil, 0, err
	}

	return units, total, nil
}

func (connection *unitConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Unit, string, bool, error) {
	var units []entity.Unit

	query := connection.db.Model(&entity.Unit{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ?", search)
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

	if err := query.Find(&units).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(units) > limit {
		last := units[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		units = units[:limit]
		hasNext = true
	}

	return units, nextCursor, hasNext, nil
}
