package repository

import (
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
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Unit, int64, error)
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

func (connection *unitConnection) Delete(id int) error {
	return connection.db.Delete(&entity.Unit{}, id).Error
}

func (connection *unitConnection) FindById(id int) (entity.Unit, error) {
	var unit entity.Unit
	err := connection.db.Preload("Products").First(&unit, id).Error
	return unit, err
}

func (connection *unitConnection) FindByBusinessId(businessId int) ([]entity.Unit, error) {
	var units []entity.Unit
	err := connection.db.Where("business_id = ?", businessId).Preload("Products").Find(&units).Error
	return units, err
}

func (connection *unitConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Unit, int64, error) {
	var units []entity.Unit
	var total int64

	// Base query dengan preload relasi
	baseQuery := connection.db.Model(&entity.Unit{}).
		Where("business_id = ?", businessId)

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
