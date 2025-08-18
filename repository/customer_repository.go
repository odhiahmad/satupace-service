package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer entity.Customer) (entity.Customer, error)
	Update(customer entity.Customer) (entity.Customer, error)
	Delete(customer entity.Customer) error
	HasRelation(customerId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(customerId uuid.UUID) (customeres entity.Customer, err error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Customer, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Customer, string, bool, error)
}

type customerConnection struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerConnection{db}
}

func (conn *customerConnection) Create(customer entity.Customer) (entity.Customer, error) {
	err := conn.db.Create(&customer).Error
	if err != nil {
		return entity.Customer{}, err
	}

	err = conn.db.First(&customer, customer.Id).Error
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}

func (conn *customerConnection) Update(customer entity.Customer) (entity.Customer, error) {
	err := conn.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&customer).Error
	if err != nil {
		return entity.Customer{}, err
	}

	err = conn.db.First(&customer, customer.Id).Error

	return customer, err
}

func (conn *customerConnection) Delete(customer entity.Customer) error {
	return conn.db.Delete(&customer).Error
}

func (conn *customerConnection) HasRelation(customerId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Product{}).Where("customer_id = ?", customerId).Count(&count).Error
	return count > 0, err
}

func (conn *customerConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Customer{}, id).Error
}

func (conn *customerConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.Customer{}, id).Error
}

func (conn *customerConnection) FindById(customerId uuid.UUID) (customeres entity.Customer, err error) {
	var customer entity.Customer
	result := conn.db.Find(&customer, customerId)
	if result != nil {
		return customer, nil
	} else {
		return customer, errors.New("tag is not found")
	}
}

func (conn *customerConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Customer, int64, error) {
	var customer []entity.Customer
	var total int64

	baseQuery := conn.db.Model(&entity.Customer{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &customer)
	if err != nil {
		return nil, 0, err
	}

	return customer, total, nil
}

func (conn *customerConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Customer, string, bool, error) {
	var customers []entity.Customer

	query := conn.db.Model(&entity.Customer{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ?", search)
	}

	sortBy := pagination.SortBy
	if sortBy == "" {
		sortBy = "updated_at"
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

	if err := query.Find(&customers).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(customers) > limit {
		last := customers[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		customers = customers[:limit]
		hasNext = true
	}

	return customers, nextCursor, hasNext, nil
}
