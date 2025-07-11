package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type ProductVariantRepository interface {
	Create(variant *entity.ProductVariant) error
	Update(variant *entity.ProductVariant) error
	Delete(id int) error
	DeleteByProductId(productId int) error
	FindById(id int) (entity.ProductVariant, error)
	FindByProductId(productId int) ([]entity.ProductVariant, error)
	SetActive(id int, active bool) error
	SetAvailable(id int, available bool) error
	IsSKUExists(sku string) (bool, error)
	CreateWithTx(txRepo ProductRepository, variants []entity.ProductVariant) error
	CountByProductId(productId int) (int64, error)
}

type ProductVariantConnection struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) ProductVariantRepository {
	return &ProductVariantConnection{db: db}
}

func (conn *ProductVariantConnection) Create(variant *entity.ProductVariant) error {
	return conn.db.Create(variant).Error
}

func (conn *ProductVariantConnection) Update(variant *entity.ProductVariant) error {
	if variant.Id == 0 {
		return errors.New("variant ID is required for update")
	}

	// Kosongkan relasi jika ada
	variant.Product = nil

	updateData := map[string]interface{}{
		"name":        variant.Name,
		"sku":         variant.SKU,
		"base_price":  variant.BasePrice,
		"stock":       variant.Stock,
		"track_stock": variant.TrackStock,
		"is_active":   variant.IsActive,
	}

	return conn.db.Model(&variant).Where("id = ?", variant.Id).Updates(updateData).Error
}

func (conn *ProductVariantConnection) Delete(id int) error {
	return conn.db.Delete(&entity.ProductVariant{}, id).Error
}

func (conn *ProductVariantConnection) DeleteByProductId(productId int) error {
	return conn.db.Where("product_id = ?", productId).Delete(&entity.ProductVariant{}).Error
}

func (conn *ProductVariantConnection) FindById(id int) (entity.ProductVariant, error) {
	var variant entity.ProductVariant
	err := conn.db.First(&variant, id).Error
	return variant, err
}

func (conn *ProductVariantConnection) FindByProductId(productId int) ([]entity.ProductVariant, error) {
	var variants []entity.ProductVariant
	err := conn.db.Where("product_id = ?", productId).Find(&variants).Error
	return variants, err
}

func (conn *ProductVariantConnection) SetActive(id int, active bool) error {
	return conn.db.Model(&entity.ProductVariant{}).
		Where("id = ?", id).
		Update("is_active", active).Error
}

func (conn *ProductVariantConnection) SetAvailable(id int, available bool) error {
	return conn.db.Model(&entity.ProductVariant{}).
		Where("id = ?", id).
		Update("is_available", available).Error
}

func (conn *ProductVariantConnection) IsSKUExists(sku string) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.ProductVariant{}).Where("sku = ?", sku).Count(&count).Error
	return count > 0, err
}

func (conn *ProductVariantConnection) CreateWithTx(txRepo ProductRepository, variants []entity.ProductVariant) error {
	tx := txRepo.(*productConnection).db
	return tx.Create(variants).Error
}

func (conn *ProductVariantConnection) CountByProductId(productId int) (int64, error) {
	var count int64
	err := conn.db.Model(&entity.ProductVariant{}).Where("product_id = ?", productId).Count(&count).Error
	return count, err
}
