package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type ProductCategoryRepository interface {
	InsertProductCategory(productCategory entity.ProductCategory) error
	UpdateProductCategory(productCategory entity.ProductCategory) error
	FindById(productCategoryId int) (entity.ProductCategory, error)
	FindAll() ([]entity.ProductCategory, error)
	FindByBusinessId(businessId int) ([]entity.ProductCategory, error)
	Delete(productCategoryId int) error
}

type ProductCategoryConnection struct {
	Db *gorm.DB
}

func NewProductCategoryRepository(Db *gorm.DB) ProductCategoryRepository {
	return &ProductCategoryConnection{Db: Db}
}

func (r *ProductCategoryConnection) InsertProductCategory(productCategory entity.ProductCategory) error {
	// Validasi keberadaan business
	var business entity.Business
	if err := r.Db.First(&business, productCategory.BusinessId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("business not found")
		}
		return err
	}

	// Optional: set default active
	productCategory.IsActive = true
	productCategory.IsAvailable = true

	// Simpan kategori
	result := r.Db.Create(&productCategory)
	return result.Error
}

func (r *ProductCategoryConnection) UpdateProductCategory(productCategory entity.ProductCategory) error {
	result := r.Db.Model(&productCategory).Updates(map[string]interface{}{
		"name":      productCategory.Name,
		"parent_id": productCategory.ParentId,
	})
	return result.Error
}

func (r *ProductCategoryConnection) FindById(productCategoryId int) (entity.ProductCategory, error) {
	var productCategory entity.ProductCategory
	result := r.Db.First(&productCategory, productCategoryId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return productCategory, errors.New("category not found")
	}
	return productCategory, result.Error
}

func (r *ProductCategoryConnection) FindAll() ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	result := r.Db.Find(&categories)
	return categories, result.Error
}

func (r *ProductCategoryConnection) FindByBusinessId(businessId int) ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	result := r.Db.Where("business_id = ?", businessId).Find(&categories)
	return categories, result.Error
}

func (r *ProductCategoryConnection) Delete(productCategoryId int) error {
	result := r.Db.Delete(&entity.ProductCategory{}, productCategoryId)
	return result.Error
}
