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

type TableRepository interface {
	Create(table entity.Table) (entity.Table, error)
	Update(table entity.Table) (entity.Table, error)
	Delete(table entity.Table) error
	HasRelation(tableId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(tableId uuid.UUID) (tablees entity.Table, err error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Table, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Table, string, bool, error)
	GetActiveTables(businessId uuid.UUID) ([]entity.Table, error)
}

type tableConnection struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) TableRepository {
	return &tableConnection{db}
}

func (conn *tableConnection) Create(table entity.Table) (entity.Table, error) {
	err := conn.db.Create(&table).Error
	if err != nil {
		return entity.Table{}, err
	}

	err = conn.db.First(&table, table.Id).Error
	if err != nil {
		return entity.Table{}, err
	}

	return table, nil
}

func (conn *tableConnection) Update(table entity.Table) (entity.Table, error) {
	err := conn.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&table).Error
	if err != nil {
		return entity.Table{}, err
	}

	err = conn.db.First(&table, table.Id).Error

	return table, err
}

func (conn *tableConnection) Delete(table entity.Table) error {
	return conn.db.Delete(&table).Error
}

func (conn *tableConnection) HasRelation(tableId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Transaction{}).Where("table_id = ?", tableId).Count(&count).Error
	return count > 0, err
}

func (conn *tableConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Table{}, id).Error
}

func (conn *tableConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.Table{}, id).Error
}

func (conn *tableConnection) FindById(tableId uuid.UUID) (tablees entity.Table, err error) {
	var table entity.Table
	result := conn.db.Find(&table, tableId)
	if result != nil {
		return table, nil
	} else {
		return table, errors.New("tag is not found")
	}
}

func (conn *tableConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Table, int64, error) {
	var table []entity.Table
	var total int64

	baseQuery := conn.db.Model(&entity.Table{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &table)
	if err != nil {
		return nil, 0, err
	}

	return table, total, nil
}

func (conn *tableConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Table, string, bool, error) {
	var tables []entity.Table

	query := conn.db.Model(&entity.Table{}).
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

	if err := query.Find(&tables).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(tables) > limit {
		last := tables[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		tables = tables[:limit]
		hasNext = true
	}

	return tables, nextCursor, hasNext, nil
}

func (conn *tableConnection) GetActiveTables(businessId uuid.UUID) ([]entity.Table, error) {
	var tables []entity.Table

	err := conn.db.
		Model(&entity.Table{}).
		Where("tables.business_id = ?", businessId).
		Preload("Business").
		Preload("Transactions", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("transactions.status IN (?) AND transactions.is_canceled = ? AND transactions.is_refunded = ?",
					[]string{"active", "in_progress"}, false, false).
				Preload("Customer").
				Preload("Cashier").
				Preload("Cashier.User").
				Preload("OrderType")
		}).
		Find(&tables).Error

	if err != nil {
		return nil, err
	}

	return tables, nil
}
