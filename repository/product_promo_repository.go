package repository

import (
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type ProductPromoRepository interface {
	CreateMany(promos []entity.ProductPromo) error
	DeleteByProductId(productId int) error
	DeleteByPromoId(promoId int) error
	FindByProductId(productId int) ([]entity.ProductPromo, error)
	FindByPromoId(promoId int) ([]entity.ProductPromo, error)
}

type productPromoRepository struct {
	db *gorm.DB
}

func NewProductPromoRepository(db *gorm.DB) ProductPromoRepository {
	return &productPromoRepository{db}
}

func (r *productPromoRepository) CreateMany(promos []entity.ProductPromo) error {
	return r.db.Create(&promos).Error
}

func (r *productPromoRepository) DeleteByProductId(productId int) error {
	return r.db.Where("product_id = ?", productId).Delete(&entity.ProductPromo{}).Error
}

func (r *productPromoRepository) DeleteByPromoId(promoId int) error {
	return r.db.Where("promo_id = ?", promoId).Delete(&entity.ProductPromo{}).Error
}

func (r *productPromoRepository) FindByProductId(productId int) ([]entity.ProductPromo, error) {
	var result []entity.ProductPromo
	err := r.db.Preload("Promo").Where("product_id = ?", productId).Find(&result).Error
	return result, err
}

func (r *productPromoRepository) FindByPromoId(promoId int) ([]entity.ProductPromo, error) {
	var result []entity.ProductPromo
	err := r.db.Preload("Product").Where("promo_id = ?", promoId).Find(&result).Error
	return result, err
}
