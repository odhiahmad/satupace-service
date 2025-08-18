package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type OrderTypeRepository interface {
	Create(orderType entity.OrderType) (entity.OrderType, error)
	Update(orderType entity.OrderType) (entity.OrderType, error)
	Delete(orderType entity.OrderType) error
	HasRelation(orderTypeId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(orderTypeId uuid.UUID) (orderTypees entity.OrderType, err error)
	FindWithPagination(pagination request.Pagination) ([]entity.OrderType, int64, error)
	FindWithPaginationCursor(pagination request.Pagination) ([]entity.OrderType, string, bool, error)
}

type orderTypeConnection struct {
	db *gorm.DB
}

func NewOrderTypeRepository(db *gorm.DB) OrderTypeRepository {
	return &orderTypeConnection{db}
}

func (conn *orderTypeConnection) Create(orderType entity.OrderType) (entity.OrderType, error) {
	err := conn.db.Create(&orderType).Error
	if err != nil {
		return entity.OrderType{}, err
	}

	err = conn.db.First(&orderType, orderType.Id).Error
	if err != nil {
		return entity.OrderType{}, err
	}

	return orderType, nil
}

func (conn *orderTypeConnection) Update(orderType entity.OrderType) (entity.OrderType, error) {
	err := conn.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&orderType).Error
	if err != nil {
		return entity.OrderType{}, err
	}

	err = conn.db.First(&orderType, orderType.Id).Error

	return orderType, err
}

func (conn *orderTypeConnection) Delete(orderType entity.OrderType) error {
	return conn.db.Delete(&orderType).Error
}

func (conn *orderTypeConnection) HasRelation(orderTypeId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Transaction{}).Where("orderType_id = ?", orderTypeId).Count(&count).Error
	return count > 0, err
}

func (conn *orderTypeConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.OrderType{}, id).Error
}

func (conn *orderTypeConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.OrderType{}, id).Error
}

func (conn *orderTypeConnection) FindById(orderTypeId uuid.UUID) (orderTypees entity.OrderType, err error) {
	var orderType entity.OrderType
	result := conn.db.Find(&orderType, orderTypeId)
	if result != nil {
		return orderType, nil
	} else {
		return orderType, errors.New("tag is not found")
	}
}

func (conn *orderTypeConnection) FindWithPagination(pagination request.Pagination) ([]entity.OrderType, int64, error) {
	var orderType []entity.OrderType
	var total int64

	baseQuery := conn.db.Model(&entity.OrderType{})

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &orderType)
	if err != nil {
		return nil, 0, err
	}

	return orderType, total, nil
}

func (conn *orderTypeConnection) FindWithPaginationCursor(pagination request.Pagination) ([]entity.OrderType, string, bool, error) {
	var orderTypes []entity.OrderType

	query := conn.db.Model(&entity.OrderType{})

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

	if err := query.Find(&orderTypes).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(orderTypes) > limit {
		last := orderTypes[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		orderTypes = orderTypes[:limit]
		hasNext = true
	}

	return orderTypes, nextCursor, hasNext, nil
}
