package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	Create(discount entity.Discount) (entity.Discount, error)
	Update(discount entity.Discount) (entity.Discount, error)
	Delete(id uuid.UUID) error
	HasRelation(brandId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	SetIsActive(id uuid.UUID, isActive bool) error
	FindById(id uuid.UUID) (entity.Discount, error)
	FindActiveGlobalDiscount(businessId uuid.UUID, now time.Time) (*entity.Discount, error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Discount, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Discount, string, bool, error)
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

func (conn *discountConnection) Delete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Discount{}, id).Error
}

func (conn *discountConnection) HasRelation(discountId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Product{}).Where("discount_id = ?", discountId).Count(&count).Error
	return count > 0, err
}

func (conn *discountConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Discount{}, id).Error
}

func (conn *discountConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.Discount{}, id).Error
}

func (conn *discountConnection) SetIsActive(id uuid.UUID, isActive bool) error {
	return conn.db.Model(&entity.Discount{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

func (conn *discountConnection) FindById(id uuid.UUID) (entity.Discount, error) {
	var discount entity.Discount
	err := conn.db.First(&discount, id).Error
	return discount, err
}

func (conn *discountConnection) FindByBusinessId(businessId uuid.UUID) ([]entity.Discount, error) {
	var discounts []entity.Discount
	err := conn.db.Where("business_id = ?", businessId).Find(&discounts).Error
	return discounts, err
}

func (conn *discountConnection) FindActiveGlobalDiscount(businessId uuid.UUID, now time.Time) (*entity.Discount, error) {
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

func (conn *discountConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Discount, int64, error) {
	var discounts []entity.Discount
	var total int64

	baseQuery := conn.db.Model(&entity.Discount{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &discounts)
	if err != nil {
		return nil, 0, err
	}

	return discounts, total, nil
}

func (conn *discountConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Discount, string, bool, error) {
	var discounts []entity.Discount

	query := conn.db.Model(&entity.Discount{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	sortBy := pagination.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "ASC"
	if pagination.OrderBy == "desc" {
		order = "DESC"
	}

	if pagination.Cursor != "" {
		cursorID, err := helper.DecodeCursorID(pagination.Cursor)
		if err != nil {
			return nil, "", false, err
		}

		if order == "ASC" {
			query = query.Where("id > ?", cursorID)
		} else {
			query = query.Where("id < ?", cursorID)
		}
	}

	limit := pagination.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, order)).Limit(limit + 1)

	if err := query.Find(&discounts).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(discounts) > limit {
		last := discounts[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		discounts = discounts[:limit]
		hasNext = true
	}

	return discounts, nextCursor, hasNext, nil
}
