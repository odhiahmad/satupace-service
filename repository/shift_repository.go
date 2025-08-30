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

type ShiftRepository interface {
	Create(shift *entity.Shift) error
	FindOpenShiftByCashier(cashierId uuid.UUID, terminalId uuid.UUID) (*entity.Shift, error)
	Update(shift *entity.Shift) error
	GetActiveShiftByTerminal(terminalId uuid.UUID) (*entity.Shift, error)
	GetActiveShiftByCashier(cashierId uuid.UUID) (*entity.Shift, error)
	FindById(shiftId uuid.UUID) (shiftes entity.Shift, err error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Shift, string, bool, error)
}

type shiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepository{db}
}

func (r *shiftRepository) Create(shift *entity.Shift) error {
	return r.db.Create(shift).Error
}

func (r *shiftRepository) FindOpenShiftByCashier(cashierId uuid.UUID, terminalId uuid.UUID) (*entity.Shift, error) {
	var shift entity.Shift
	err := r.db.
		Preload("Business").
		Preload("Cashier").
		Preload("Terminal").
		Where("cashier_id = ? AND terminal_id = ? AND status = ?", cashierId, terminalId, "open").First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepository) Update(shift *entity.Shift) error {
	return r.db.Save(shift).Error
}

func (r *shiftRepository) GetActiveShiftByTerminal(terminalId uuid.UUID) (*entity.Shift, error) {
	var shift entity.Shift
	err := r.db.
		Preload("Business").
		Preload("Cashier").
		Preload("Terminal").
		Where("terminal_id = ? AND status = ? AND closed_at IS NULL", terminalId, "open").
		First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepository) GetActiveShiftByCashier(cashierId uuid.UUID) (*entity.Shift, error) {
	var shift entity.Shift
	err := r.db.
		Preload("Business").
		Preload("Cashier").
		Preload("Terminal").
		Where("cashier_id = ? AND status = ? AND closed_at IS NULL", cashierId, "open").
		First(&shift).Error

	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepository) FindById(shiftId uuid.UUID) (shiftes entity.Shift, err error) {
	var shift entity.Shift
	result := r.db.
		Preload("Business").
		Preload("Cashier").
		Preload("Terminal").
		Find(&shift, shiftId)
	if result != nil {
		return shift, nil
	} else {
		return shift, errors.New("tag is not found")
	}
}

func (r *shiftRepository) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Shift, string, bool, error) {
	var shifts []entity.Shift

	query := r.db.Model(&entity.Shift{}).
		Preload("Business").
		Preload("Cashier").
		Preload("Terminal").
		Where("business_id = ?", businessId)

	if pagination.Status != "" {
		query = query.Where("status = ?", pagination.Status)
	}

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Joins("JOIN user_businesses u ON u.id = shifts.cashier_id").
			Where("u.name ILIKE ?", search)
	}

	sortBy := pagination.SortBy
	if sortBy == "" {
		sortBy = "opened_at"
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

	if err := query.Find(&shifts).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(shifts) > limit {
		last := shifts[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		shifts = shifts[:limit]
		hasNext = true
	}

	return shifts, nextCursor, hasNext, nil
}
