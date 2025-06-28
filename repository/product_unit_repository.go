package repository

import (
	"time"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type UnitRepository interface {
	Create(unit entity.Unit) (entity.Unit, error)
	Update(unit entity.Unit) (entity.Unit, error)
	Delete(id int) error
	FindById(id int) (entity.Unit, error)
	FindActiveGlobalUnit(businessId int, now time.Time) (*entity.Unit, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Unit, int64, error)
}

type unitRepo struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepo{db: db}
}

func (r *unitRepo) Create(unit entity.Unit) (entity.Unit, error) {
	err := r.db.Create(&unit).Error
	return unit, err
}

func (r *unitRepo) Update(unit entity.Unit) (entity.Unit, error) {
	err := r.db.Save(&unit).Error
	return unit, err
}

func (r *unitRepo) Delete(id int) error {
	return r.db.Delete(&entity.Unit{}, id).Error
}

func (r *unitRepo) FindById(id int) (entity.Unit, error) {
	var unit entity.Unit
	err := r.db.Preload("Products").First(&unit, id).Error
	return unit, err
}

func (r *unitRepo) FindByBusinessId(businessId int) ([]entity.Unit, error) {
	var units []entity.Unit
	err := r.db.Where("business_id = ?", businessId).Preload("Products").Find(&units).Error
	return units, err
}

func (r *unitRepo) FindActiveGlobalUnit(businessId int, now time.Time) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.
		Where("business_id = ? AND is_global = ? AND start_at <= ? AND end_at >= ?", businessId, true, now, now).
		Order("start_at DESC").
		First(&unit).Error

	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepo) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Unit, int64, error) {
	var units []entity.Unit
	var total int64

	// Base query dengan preload relasi
	baseQuery := r.db.Model(&entity.Unit{}).
		Where("business_id = ?", businessId).
		Preload("Products")

	// Filter pencarian
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Hitung total data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &units)
	if err != nil {
		return nil, 0, err
	}

	return units, total, nil
}
