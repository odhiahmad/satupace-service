package repository

import (
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type ShiftRepository interface {
	Create(shift *entity.Shift) error
	FindOpenShiftByCashier(cashierId uuid.UUID, terminalId uuid.UUID) (*entity.Shift, error)
	Update(shift *entity.Shift) error
	GetActiveShiftByTerminal(terminalId uuid.UUID) (*entity.Shift, error)
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
	err := r.db.Where("cashier_id = ? AND terminal_id = ? AND status = ?", cashierId, terminalId, "open").First(&shift).Error
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
		Preload("Cashier").
		Where("terminal_id = ? AND status = ? AND closed_at IS NULL", terminalId, "open").
		First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}
