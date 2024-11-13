package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type ProductUnitRepository interface {
	InsertProductUnit(productUnit entity.ProductUnit)
	UpdateProductUnit(productUnit entity.ProductUnit)
	FindById(productUnitId int) (productUnit entity.ProductUnit, err error)
	FindAll() []entity.ProductUnit
	Delete(productUnitId int)
}

type ProductUnitConnection struct {
	Db *gorm.DB
}

func NewProductUnitRepository(Db *gorm.DB) ProductUnitRepository {
	return &ProductUnitConnection{Db: Db}
}

func (t *ProductUnitConnection) InsertProductUnit(productUnit entity.ProductUnit) {
	result := t.Db.Create(&productUnit)

	helper.ErrorPanic(result.Error)
}

func (t *ProductUnitConnection) UpdateProductUnit(productUnit entity.ProductUnit) {
	var updateProductUnit = request.ProductUnitUpdate{
		Id:   productUnit.Id,
		Name: productUnit.Name,
	}

	result := t.Db.Model(&productUnit).Updates(updateProductUnit)
	helper.ErrorPanic(result.Error)
}

func (t *ProductUnitConnection) FindById(productUnitId int) (productUnits entity.ProductUnit, err error) {
	var productUnit entity.ProductUnit
	result := t.Db.Find(&productUnit, productUnitId)
	if result != nil {
		return productUnit, nil
	} else {
		return productUnit, errors.New("tag is not found")
	}
}

func (t *ProductUnitConnection) FindAll() []entity.ProductUnit {
	var productUnit []entity.ProductUnit
	result := t.Db.Find(&productUnit)
	helper.ErrorPanic(result.Error)
	return productUnit
}

func (t *ProductUnitConnection) Delete(productUnitId int) {
	var productUnits entity.ProductUnit
	result := t.Db.Where("id = ?", productUnitId).Delete(&productUnits)
	helper.ErrorPanic(result.Error)
}
