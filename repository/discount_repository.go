package repository

import (
	"time"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	Create(discount entity.Discount) (entity.Discount, error)
	Update(discount entity.Discount) (entity.Discount, error)
	Delete(id int) error
	FindById(id int) (entity.Discount, error)
	FindActiveGlobalDiscount(businessId int, now time.Time) (*entity.Discount, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Discount, int64, error)
}

type discountRepo struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepo{db: db}
}

func (r *discountRepo) Create(discount entity.Discount) (entity.Discount, error) {
	err := r.db.Create(&discount).Error
	return discount, err
}

func (r *discountRepo) Update(discount entity.Discount) (entity.Discount, error) {
	err := r.db.Save(&discount).Error
	return discount, err
}

func (r *discountRepo) Delete(id int) error {
	return r.db.Delete(&entity.Discount{}, id).Error
}

func (r *discountRepo) FindById(id int) (entity.Discount, error) {
	var discount entity.Discount
	err := r.db.Preload("Products").First(&discount, id).Error
	return discount, err
}

func (r *discountRepo) FindByBusinessId(businessId int) ([]entity.Discount, error) {
	var discounts []entity.Discount
	err := r.db.Where("business_id = ?", businessId).Preload("Products").Find(&discounts).Error
	return discounts, err
}

func (r *discountRepo) FindActiveGlobalDiscount(businessId int, now time.Time) (*entity.Discount, error) {
	var discount entity.Discount
	err := r.db.
		Where("business_id = ? AND is_global = ? AND start_at <= ? AND end_at >= ?", businessId, true, now, now).
		Order("start_at DESC").
		First(&discount).Error

	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *discountRepo) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Discount, int64, error) {
	var discounts []entity.Discount
	var total int64

	// Base query
	baseQuery := r.db.Model(&entity.Discount{}).
		Where("business_id = ?", businessId).
		Preload("Products")

	// Search
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Hitung total data sebelum pagination
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Paginate query
	_, _, err := p.Paginate(baseQuery, &discounts)
	if err != nil {
		return nil, 0, err
	}

	return discounts, total, nil
}
