package repository

import (
	"errors"

	"loka-kasir/entity"
	"loka-kasir/helper"

	"gorm.io/gorm"
)

type BusinessTypeRepository interface {
	InsertBusinessType(businessType entity.BusinessType)
	UpdateBusinessType(businessType entity.BusinessType)
	FindById(businessTypeId int) (businessType entity.BusinessType, err error)
	FindAll() []entity.BusinessType
	Delete(businessTypeId int)
}

type businessTypeConnection struct {
	db *gorm.DB
}

func NewBusinessTypeRepository(db *gorm.DB) BusinessTypeRepository {
	return &businessTypeConnection{db: db}
}

func (conn *businessTypeConnection) InsertBusinessType(businessType entity.BusinessType) {
	result := conn.db.Create(&businessType)

	helper.ErrorPanic(result.Error)
}

func (conn *businessTypeConnection) UpdateBusinessType(businessType entity.BusinessType) {
	// Pastikan ID tidak kosong
	if businessType.Id == 0 {
		helper.ErrorPanic(errors.New("ID businessType tidak boleh kosong"))
	}

	// Update langsung menggunakan struct entity
	result := conn.db.Model(&entity.BusinessType{}).
		Where("id = ?", businessType.Id).
		Updates(map[string]interface{}{
			"name": businessType.Name,
		})

	helper.ErrorPanic(result.Error)
}

func (conn *businessTypeConnection) FindById(businessTypeId int) (businessTypes entity.BusinessType, err error) {
	var businessType entity.BusinessType
	result := conn.db.Find(&businessType, businessTypeId)
	if result != nil {
		return businessType, nil
	} else {
		return businessType, errors.New("tag is not found")
	}
}

func (conn *businessTypeConnection) FindAll() []entity.BusinessType {
	var businessType []entity.BusinessType
	result := conn.db.Find(&businessType)
	helper.ErrorPanic(result.Error)
	return businessType
}

func (conn *businessTypeConnection) Delete(businessTypeId int) {
	var businessTypes entity.BusinessType
	result := conn.db.Where("id = ?", businessTypeId).Delete(&businessTypes)
	helper.ErrorPanic(result.Error)
}
