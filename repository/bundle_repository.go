package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type BundleRepository interface {
	InsertBundle(bundle *entity.Bundle) error
	UpdateBundle(bundle *entity.Bundle) error
	FindById(bundleId int) (entity.Bundle, error)
	FindAll() ([]entity.Bundle, error)
	Delete(bundleId int) error
	InsertItemsByBundleId(bundleId int, items []entity.BundleItem) error
	DeleteItemsByBundleId(bundleId int) error
	FindByBusinessId(businessId int) ([]entity.Bundle, error)
}

type BundleConnection struct {
	Db *gorm.DB
}

func NewBundleRepository(Db *gorm.DB) BundleRepository {
	return &BundleConnection{Db: Db}
}

func (r *BundleConnection) InsertBundle(bundle *entity.Bundle) error {
	result := r.Db.Create(bundle)
	return result.Error
}

func (r *BundleConnection) UpdateBundle(bundle *entity.Bundle) error {
	result := r.Db.Save(bundle)
	return result.Error
}

func (r *BundleConnection) FindById(bundleId int) (entity.Bundle, error) {
	var bundle entity.Bundle
	result := r.Db.Preload("Items.Product").First(&bundle, bundleId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return bundle, errors.New("bundle product not found")
	}
	return bundle, result.Error
}

func (r *BundleConnection) FindAll() ([]entity.Bundle, error) {
	var bundles []entity.Bundle
	result := r.Db.Preload("Items.Product").Find(&bundles)
	return bundles, result.Error
}

func (r *BundleConnection) Delete(bundleId int) error {
	if err := r.DeleteItemsByBundleId(bundleId); err != nil {
		return err
	}
	result := r.Db.Delete(&entity.Bundle{}, bundleId)
	return result.Error
}

func (r *BundleConnection) InsertItemsByBundleId(bundleId int, items []entity.BundleItem) error {
	for i := range items {
		items[i].BundleId = bundleId
	}
	result := r.Db.Create(&items)
	return result.Error
}

func (r *BundleConnection) DeleteItemsByBundleId(bundleId int) error {
	result := r.Db.Where("bundle_id = ?", bundleId).Delete(&entity.BundleItem{})
	return result.Error
}

func (r *BundleConnection) FindByBusinessId(businessId int) ([]entity.Bundle, error) {
	var bundles []entity.Bundle
	result := r.Db.Preload("Items.Product").Where("business_id = ?", businessId).Find(&bundles)
	return bundles, result.Error
}
