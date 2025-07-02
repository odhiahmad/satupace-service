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
	AppendRequiredProducts(promo *entity.Promo, products []entity.Product) error
	ReplaceRequiredProducts(promoId int, products []entity.Product) error
	SetIsActive(id int, isActive bool) error
	Exists(id int) (bool, error)
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db}
}

func (r *promoRepository) Create(promo entity.Promo) (entity.Promo, error) {
	err := r.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&promo).Error
	return promo, err
}

func (r *promoRepository) Update(promo entity.Promo) (entity.Promo, error) {
	err := r.db.Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Select("*").Updates(&promo).Error
	return promo, err
}

func (r *promoRepository) Delete(promo entity.Promo) error {
	// Hapus relasi required_products juga (opsional, untuk bersih)
	_ = r.db.Model(&promo).Association("RequiredProducts").Clear()
	return r.db.Delete(&promo).Error
}

func (r *promoRepository) AppendRequiredProducts(promo *entity.Promo, products []entity.Product) error {
	return r.db.Model(promo).Association("RequiredProducts").Append(&products)
}

func (r *promoRepository) FindById(id int) (entity.Promo, error) {
	var promo entity.Promo
	err := r.db.
		Preload("RequiredProducts"). // ⬅️ Tambahkan preload untuk required
		First(&promo, id).Error
	return promo, err
}

func (r *promoRepository) ReplaceRequiredProducts(promoId int, products []entity.Product) error {
	promo := entity.Promo{Id: promoId}
	return r.db.Model(&promo).Association("RequiredProducts").Replace(products)
}

func (r *promoRepository) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Promo, int64, error) {
	var promos []entity.Promo
	var total int64

	baseQuery := r.db.Model(&entity.Promo{}).
		Preload("RequiredProducts").
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination)
	_, _, err := p.Paginate(baseQuery, &promos)
	if err != nil {
		return nil, 0, err
	}

	return promos, total, nil
}

func (r *promoRepository) SetIsActive(id int, active bool) error {
	return r.db.Model(&entity.Promo{}).Where("id = ?", id).Update("is_active", active).Error
}

func (r *promoRepository) Exists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Promo{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
