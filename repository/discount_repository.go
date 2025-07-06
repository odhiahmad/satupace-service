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
	SetIsActive(id int, isActive bool) error
}

type discountConnection struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountConnection{db: db}
}

func (conn *discountConnection) Create(discount entity.Discount) (entity.Discount, error) {
	err := conn.db.Create(&discount).Error
	return discount, err
}

func (conn *discountConnection) Update(discount entity.Discount) (entity.Discount, error) {
	err := conn.db.Model(&entity.Discount{}).Where("id = ?", discount.Id).Updates(discount).Error
	if err != nil {
		return entity.Discount{}, err
	}
	return discount, nil
}

func (conn *discountConnection) Delete(id int) error {
	return conn.db.Delete(&entity.Discount{}, id).Error
}

func (conn *discountConnection) FindById(id int) (entity.Discount, error) {
	var discount entity.Discount
	err := conn.db.First(&discount, id).Error
	return discount, err
}

func (conn *discountConnection) FindByBusinessId(businessId int) ([]entity.Discount, error) {
	var discounts []entity.Discount
	err := conn.db.Where("business_id = ?", businessId).Find(&discounts).Error
	return discounts, err
}

func (conn *discountConnection) FindActiveGlobalDiscount(businessId int, now time.Time) (*entity.Discount, error) {
	var discount entity.Discount
	err := conn.db.
		Where("business_id = ? AND is_global = ? AND start_at <= ? AND end_at >= ?", businessId, true, now, now).
		Order("start_at DESC").
		First(&discount).Error

	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (conn *discountConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Discount, int64, error) {
	var discounts []entity.Discount
	var total int64

	// Base query
	baseQuery := conn.db.Model(&entity.Discount{}).
		Where("business_id = ?", businessId)

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
	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	// Paginate query
	_, _, err := p.Paginate(baseQuery, &discounts)
	if err != nil {
		return nil, 0, err
	}

	return discounts, total, nil
}

func (conn *discountConnection) SetIsActive(id int, isActive bool) error {
	return conn.db.Model(&entity.Discount{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}
