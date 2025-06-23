package repository

import (
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type PromoRepository interface {
	Create(promo entity.Promo) (entity.Promo, error)
	Update(promo entity.Promo) (entity.Promo, error)
	Delete(promo entity.Promo) error
	FindById(id int) (entity.Promo, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Promo, int64, error)
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db}
}

func (r *promoRepository) Create(promo entity.Promo) (entity.Promo, error) {
	err := r.db.Create(&promo).Error
	return promo, err
}

func (r *promoRepository) Update(promo entity.Promo) (entity.Promo, error) {
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&promo).Error
	return promo, err
}

func (r *promoRepository) Delete(promo entity.Promo) error {
	return r.db.Delete(&promo).Error
}

func (r *promoRepository) FindById(id int) (entity.Promo, error) {
	var promo entity.Promo
	err := r.db.Preload("ProductPromos.Product").First(&promo, id).Error
	return promo, err
}

func (r *promoRepository) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Promo, int64, error) {
	var promos []entity.Promo
	var total int64

	// Base query dengan preload relasi
	baseQuery := r.db.Model(&entity.Promo{}).
		Where("business_id = ?", businessId).
		Preload("ProductPromos.Product")

	// Filter pencarian
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Hitung total data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &promos)
	if err != nil {
		return nil, 0, err
	}

	return promos, total, nil
}
