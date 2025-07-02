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

type productCategoryConnection struct {
	Db *gorm.DB
}

func NewProductCategoryRepository(Db *gorm.DB) ProductCategoryRepository {
	return &productCategoryConnection{Db: Db}
}

func (conn *productCategoryConnection) InsertProductCategory(productCategory entity.ProductCategory) error {
	// Validasi keberadaan business
	var business entity.Business
	if err := conn.Db.First(&business, productCategory.BusinessId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("business not found")
		}
		return err
	}

	productCategory.IsActive = true
	result := conn.Db.Create(&productCategory)
	return result.Error
}

func (conn *productCategoryConnection) UpdateProductCategory(productCategory entity.ProductCategory) error {
	result := conn.Db.Model(&productCategory).Updates(map[string]interface{}{
		"name":      productCategory.Name,
		"parent_id": productCategory.ParentId,
	})
	return result.Error
}

func (conn *productCategoryConnection) FindById(productCategoryId int) (entity.ProductCategory, error) {
	var productCategory entity.ProductCategory
	result := conn.Db.First(&productCategory, productCategoryId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return productCategory, errors.New("category not found")
	}
	return productCategory, result.Error
}

func (conn *productCategoryConnection) FindAll() ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	result := conn.Db.Find(&categories)
	return categories, result.Error
}

func (conn *productCategoryConnection) FindByBusinessId(businessId int) ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	result := conn.Db.Where("business_id = ?", businessId).Find(&categories)
	return categories, result.Error
}

func (conn *productCategoryConnection) Delete(productCategoryId int) error {
	result := conn.Db.Delete(&entity.ProductCategory{}, productCategoryId)
	return result.Error
}
