package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BusinessTypeRepository interface {
	InsertBusinessType(businessType entity.BusinessType)
	UpdateBusinessType(businessType entity.BusinessType)
	FindById(businessTypeId int) (businessType entity.BusinessType, err error)
	FindAll() []entity.BusinessType
	Delete(businessTypeId int)
}

type BusinessTypeConnection struct {
	Db *gorm.DB
}

func NewBusinessTypeRepository(Db *gorm.DB) BusinessTypeRepository {
	return &BusinessTypeConnection{Db: Db}
}

func (t *BusinessTypeConnection) InsertBusinessType(businessType entity.BusinessType) {
	result := t.Db.Create(&businessType)

	helper.ErrorPanic(result.Error)
}

func (t *BusinessTypeConnection) UpdateBusinessType(businessType entity.BusinessType) {
	var updateBusinessType = request.BusinessTypeUpdate{
		Id:   businessType.Id,
		Name: businessType.Name,
	}

	result := t.Db.Model(&businessType).Updates(updateBusinessType)
	helper.ErrorPanic(result.Error)
}

func (t *BusinessTypeConnection) FindById(businessTypeId int) (businessTypes entity.BusinessType, err error) {
	var businessType entity.BusinessType
	result := t.Db.Find(&businessType, businessTypeId)
	if result != nil {
		return businessType, nil
	} else {
		return businessType, errors.New("tag is not found")
	}
}

func (t *BusinessTypeConnection) FindAll() []entity.BusinessType {
	var businessType []entity.BusinessType
	result := t.Db.Find(&businessType)
	helper.ErrorPanic(result.Error)
	return businessType
}

func (t *BusinessTypeConnection) Delete(businessTypeId int) {
	var businessTypes entity.BusinessType
	result := t.Db.Where("id = ?", businessTypeId).Delete(&businessTypes)
	helper.ErrorPanic(result.Error)
}
