package repository

import (
	"errors"
	"fmt"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type TaxRepository interface {
	Create(tax entity.Tax) (entity.Tax, error)
	Update(tax entity.Tax) (entity.Tax, error)
	Delete(tax entity.Tax) error
	FindById(taxId int) (taxes entity.Tax, err error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Tax, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]entity.Tax, string, error)
}

type taxConnection struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) TaxRepository {
	return &taxConnection{db}
}

func (conn *taxConnection) Create(tax entity.Tax) (entity.Tax, error) {
	err := conn.db.Create(&tax).Error
	if err != nil {
		return entity.Tax{}, err
	}

	err = conn.db.First(&tax, tax.Id).Error
	if err != nil {
		return entity.Tax{}, err
	}

	return tax, nil
}

func (conn *taxConnection) Update(tax entity.Tax) (entity.Tax, error) {
	err := conn.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&tax).Error
	if err != nil {
		return entity.Tax{}, err
	}

	err = conn.db.First(&tax, tax.Id).Error

	return tax, err
}

func (conn *taxConnection) Delete(tax entity.Tax) error {
	return conn.db.Delete(&tax).Error
}

func (conn *taxConnection) FindById(taxId int) (taxes entity.Tax, err error) {
	var tax entity.Tax
	result := conn.db.Find(&tax, taxId)
	if result != nil {
		return tax, nil
	} else {
		return tax, errors.New("tag is not found")
	}
}

func (conn *taxConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Tax, int64, error) {
	var tax []entity.Tax
	var total int64

	baseQuery := conn.db.Model(&entity.Tax{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &tax)
	if err != nil {
		return nil, 0, err
	}

	return tax, total, nil
}

func (conn *taxConnection) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]entity.Tax, string, error) {
	var taxes []entity.Tax

	query := conn.db.Model(&entity.Tax{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ?", search)
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
			return nil, "", err
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

	if err := query.Find(&taxes).Error; err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(taxes) > limit {
		last := taxes[limit-1]
		nextCursor = helper.EncodeCursorID(int64(last.Id))
		taxes = taxes[:limit]
	}

	return taxes, nextCursor, nil
}
