package repository

import (
	"fmt"

	"loka-kasir/data/request"
	"loka-kasir/entity"
	"loka-kasir/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(employee entity.Employee) (entity.Employee, error)
	Update(employee entity.Employee) (entity.Employee, error)
	Delete(id uuid.UUID) error
	HasRelation(brandId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(id uuid.UUID) (entity.Employee, error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Employee, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Employee, string, bool, error)
}

type employeeConnection struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeConnection{db: db}
}

func (connection *employeeConnection) Create(employee entity.Employee) (entity.Employee, error) {
	err := connection.db.Create(&employee).Error
	return employee, err
}

func (connection *employeeConnection) Update(employee entity.Employee) (entity.Employee, error) {
	err := connection.db.Save(&employee).Error
	return employee, err
}

func (connection *employeeConnection) Delete(id uuid.UUID) error {
	return connection.db.Delete(&entity.Employee{}, id).Error
}

func (conn *employeeConnection) HasRelation(employeeId uuid.UUID) (bool, error) {
	var count int64
	err := conn.db.Model(&entity.Product{}).Where("employee_id = ?", employeeId).Count(&count).Error
	return count > 0, err
}

func (conn *employeeConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.Employee{}, id).Error
}

func (conn *employeeConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.Employee{}, id).Error
}

func (connection *employeeConnection) FindById(id uuid.UUID) (entity.Employee, error) {
	var employee entity.Employee
	err := connection.db.First(&employee, id).Error
	return employee, err
}

func (connection *employeeConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Employee, int64, error) {
	var employees []entity.Employee
	var total int64

	baseQuery := connection.db.Model(&entity.Employee{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &employees)
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (connection *employeeConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Employee, string, bool, error) {
	var employees []entity.Employee

	query := connection.db.Model(&entity.Employee{}).
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

	if err := query.Find(&employees).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(employees) > limit {
		last := employees[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		employees = employees[:limit]
		hasNext = true
	}

	return employees, nextCursor, hasNext, nil
}
