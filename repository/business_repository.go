package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BusinessRepository interface {
	InsertBusiness(business entity.Business) (entity.Business, error)
	FindById(businessId int) (business entity.Business, err error)
	FindAll() []entity.Business
	Delete(businessId int)
}

type BusinessConnection struct {
	Db *gorm.DB
}

func NewBusinessRepository(Db *gorm.DB) BusinessRepository {
	return &BusinessConnection{Db: Db}
}

func (t *BusinessConnection) InsertBusiness(business entity.Business) (entity.Business, error) {
	result := t.Db.Create(&business)
	helper.ErrorPanic(result.Error)

	return business, result.Error
}

func (t *BusinessConnection) FindById(businessId int) (businesss entity.Business, err error) {
	var business entity.Business
	result := t.Db.First(&business, businessId)
	if result.Error != nil {
		return business, errors.New("business not found")
	}
	return business, nil
}

func (t *BusinessConnection) FindAll() []entity.Business {
	var business []entity.Business
	result := t.Db.Find(&business)
	helper.ErrorPanic(result.Error)
	return business
}

func (t *BusinessConnection) Delete(businessId int) {
	var businesss entity.Business
	result := t.Db.Where("id = ?", businessId).Delete(&businesss)
	helper.ErrorPanic(result.Error)
}
