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

type TerminalRepository interface {
	Create(terminal entity.Terminal) (entity.Terminal, error)
	Update(terminal entity.Terminal) (entity.Terminal, error)
	Delete(terminal entity.Terminal) error
	HasRelation(terminalId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(terminalId uuid.UUID) (terminales entity.Terminal, err error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Terminal, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Terminal, string, bool, error)
}

type terminalConnection struct {
	db *gorm.DB
}

func NewTerminalRepository(db *gorm.DB) TerminalRepository {
	return &terminalConnection{db}
}

func (conn *terminalConnection) Create(terminal entity.Terminal) (entity.Terminal, error) {
	err := conn.db.Create(&terminal).Error
	if err != nil {
		return entity.Terminal{}, err
	}

	err = conn.db.First(&terminal, terminal.Id).Error
	if err != nil {
		return entity.Terminal{}, err
	}

	return terminal, nil
}

func (conn *terminalConnection) Update(terminal entity.Terminal) (entity.Terminal, error) {
	err := conn.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&terminal).Error
	if err != nil {
		return entity.Terminal{}, err
	}

	err = conn.db.First(&terminal, terminal.Id).Error

	return terminal, err
}

func (conn *terminalConnection) Delete(terminal entity.Terminal) error {
	return conn.db.Delete(&terminal).Error
}

func (conn *terminalConnection) HasRelation(terminalId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Product{}).Where("terminal_id = ?", terminalId).Count(&count).Error
	return count > 0, err
}

func (conn *terminalConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Terminal{}, id).Error
}

func (conn *terminalConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.Terminal{}, id).Error
}

func (conn *terminalConnection) FindById(terminalId uuid.UUID) (terminales entity.Terminal, err error) {
	var terminal entity.Terminal
	result := conn.db.Find(&terminal, terminalId)
	if result != nil {
		return terminal, nil
	} else {
		return terminal, errors.New("tag is not found")
	}
}

func (conn *terminalConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Terminal, int64, error) {
	var terminal []entity.Terminal
	var total int64

	baseQuery := conn.db.Model(&entity.Terminal{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &terminal)
	if err != nil {
		return nil, 0, err
	}

	return terminal, total, nil
}

func (conn *terminalConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Terminal, string, bool, error) {
	var terminals []entity.Terminal

	query := conn.db.Model(&entity.Terminal{}).
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

	if err := query.Find(&terminals).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(terminals) > limit {
		last := terminals[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		terminals = terminals[:limit]
		hasNext = true
	}

	return terminals, nextCursor, hasNext, nil
}
