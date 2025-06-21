package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BundleRepository interface {
	InsertBundle(bundle entity.Bundle) entity.Bundle
	UpdateBundle(bundle entity.Bundle)
	FindById(bundleId int) (entity.Bundle, error)
	FindAll() []entity.Bundle
	Delete(bundleId int)
	InsertItemsByBundleId(bundleId int, items []entity.BundleItem)
	DeleteItemsByBundleId(bundleId int)
}

type BundleConnection struct {
	Db *gorm.DB
}

func NewBundleRepository(Db *gorm.DB) BundleRepository {
	return &BundleConnection{Db: Db}
}

func (r *BundleConnection) InsertBundle(bundle entity.Bundle) entity.Bundle {
	result := r.Db.Create(&bundle)
	helper.ErrorPanic(result.Error)
	return bundle
}

func (r *BundleConnection) UpdateBundle(bundle entity.Bundle) {
	result := r.Db.Save(&bundle)
	helper.ErrorPanic(result.Error)
}

func (r *BundleConnection) FindById(bundleId int) (entity.Bundle, error) {
	var bundle entity.Bundle
	result := r.Db.Preload("Items.Product").First(&bundle, bundleId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return bundle, errors.New("bundle product not found")
	}
	return bundle, result.Error
}

func (r *BundleConnection) FindAll() []entity.Bundle {
	var bundles []entity.Bundle
	result := r.Db.Preload("Items.Product").Find(&bundles)
	helper.ErrorPanic(result.Error)
	return bundles
}

func (r *BundleConnection) Delete(bundleId int) {
	r.DeleteItemsByBundleId(bundleId)
	result := r.Db.Delete(&entity.Bundle{}, bundleId)
	helper.ErrorPanic(result.Error)
}

func (r *BundleConnection) InsertItemsByBundleId(bundleId int, items []entity.BundleItem) {
	for i := range items {
		items[i].Id = bundleId
	}
	result := r.Db.Create(&items)
	helper.ErrorPanic(result.Error)
}

func (r *BundleConnection) DeleteItemsByBundleId(bundleId int) {
	result := r.Db.Where("bundle_product_id = ?", bundleId).Delete(&entity.Bundle{})
	helper.ErrorPanic(result.Error)
}
