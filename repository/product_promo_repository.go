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
	CreateManyWithTx(txRepo ProductRepository, promos []entity.ProductPromo) error
}

type productPromoConnection struct {
	db *gorm.DB
}

func NewProductPromoRepository(db *gorm.DB) ProductPromoRepository {
	return &productPromoConnection{db}
}

func (conn *productPromoConnection) CreateMany(promos []entity.ProductPromo) error {
	return conn.db.Create(&promos).Error
}

func (conn *productPromoConnection) DeleteByProductId(productId int) error {
	return conn.db.Where("product_id = ?", productId).Delete(&entity.ProductPromo{}).Error
}

func (conn *productPromoConnection) DeleteByPromoId(promoId int) error {
	return conn.db.Where("promo_id = ?", promoId).Delete(&entity.ProductPromo{}).Error
}

func (conn *productPromoConnection) FindByProductId(productId int) ([]entity.ProductPromo, error) {
	var result []entity.ProductPromo
	err := conn.db.Preload("Promo").Where("product_id = ?", productId).Find(&result).Error
	return result, err
}

func (conn *productPromoConnection) FindByPromoId(promoId int) ([]entity.ProductPromo, error) {
	var result []entity.ProductPromo
	err := conn.db.Preload("Product").Where("promo_id = ?", promoId).Find(&result).Error
	return result, err
}

func (conn *productPromoConnection) CreateManyWithTx(txRepo ProductRepository, promos []entity.ProductPromo) error {
	tx := txRepo.(*productConnection).DB()
	return tx.Create(&promos).Error
}
