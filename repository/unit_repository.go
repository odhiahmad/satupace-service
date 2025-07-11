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
	err := connection.db.First(&unit, id).Error
	return unit, err
}

func (connection *unitConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Unit, int64, error) {
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
