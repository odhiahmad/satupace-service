package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	InsertProduct(product *entity.Product) error
	UpdateProduct(product *entity.Product) error
	FindById(productId int) (entity.Product, error)
	FindAll() ([]entity.Product, error)
	Delete(productId int) error
}

type ProductConnection struct {
	Db *gorm.DB
}

func NewProductRepository(Db *gorm.DB) ProductRepository {
	return &ProductConnection{Db: Db}
}

func (r *ProductConnection) InsertProduct(product *entity.Product) error {
	result := r.Db.Create(&product)
	return result.Error
}

func (r *ProductConnection) UpdateProduct(product *entity.Product) error {
	result := r.Db.Save(&product)
	return result.Error
}

func (r *ProductConnection) FindById(productId int) (entity.Product, error) {
	var product entity.Product
	result := r.Db.Preload("Variants").Preload("Business").Preload("ProductCategory").First(&product, productId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return product, errors.New("product not found")
	}
	return product, result.Error
}

func (r *ProductConnection) FindAll() ([]entity.Product, error) {
	var products []entity.Product
	result := r.Db.Preload("Variants").Preload("Business").Preload("ProductCategory").Find(&products)
	return products, result.Error
}

func (r *ProductConnection) Delete(productId int) error {
	result := r.Db.Delete(&entity.Product{}, productId)
	return result.Error
}
