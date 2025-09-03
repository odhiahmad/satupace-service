package repository

import (
	"errors"
	"fmt"

	"loka-kasir/data/request"
	"loka-kasir/entity"
	"loka-kasir/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	InsertCategory(category entity.Category) (entity.Category, error)
	UpdateCategory(category entity.Category) (entity.Category, error)
	Delete(categoryId uuid.UUID) error
	HasRelation(categoryId uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindById(categoryId uuid.UUID) (entity.Category, error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Category, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Category, string, bool, error)
	FindWithPaginationCursorProduct(businessId uuid.UUID, pagination request.Pagination) ([]entity.Category, string, bool, error)
}

type categoryConnection struct {
	Db *gorm.DB
}

func NewCategoryRepository(Db *gorm.DB) CategoryRepository {
	return &categoryConnection{Db: Db}
}

func (conn *categoryConnection) InsertCategory(category entity.Category) (entity.Category, error) {
	var business entity.Business
	if err := conn.Db.First(&business, category.BusinessId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Category{}, errors.New("business not found")
		}
		return entity.Category{}, err
	}

	category.IsActive = true
	err := conn.Db.Create(&category).Error
	if err != nil {
		return entity.Category{}, err
	}

	err = conn.Db.First(&category, category.Id).Error
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (conn *categoryConnection) UpdateCategory(category entity.Category) (entity.Category, error) {
	err := conn.Db.Model(&category).Updates(map[string]interface{}{
		"name":      category.Name,
		"parent_id": category.ParentId,
	}).Error
	if err != nil {
		return entity.Category{}, err
	}

	err = conn.Db.First(&category, category.Id).Error
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (conn *categoryConnection) Delete(categoryId uuid.UUID) error {
	result := conn.Db.Delete(&entity.Category{}, categoryId)
	return result.Error
}

func (conn *categoryConnection) HasRelation(categoryId uuid.UUID) (bool, error) {
	var count int64
	err := conn.Db.Model(&entity.Product{}).Where("category_id = ?", categoryId).Count(&count).Error
	return count > 0, err
}

func (conn *categoryConnection) SoftDelete(id uuid.UUID) error {
	return conn.Db.Delete(&entity.Category{}, id).Error
}

func (conn *categoryConnection) HardDelete(id uuid.UUID) error {
	return conn.Db.Unscoped().Delete(&entity.Category{}, id).Error
}

func (conn *categoryConnection) FindById(categoryId uuid.UUID) (entity.Category, error) {
	var category entity.Category
	result := conn.Db.First(&category, categoryId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return category, errors.New("category not found")
	}
	return category, result.Error
}

func (conn *categoryConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Category, int64, error) {
	var category []entity.Category
	var total int64

	baseQuery := conn.Db.Model(&entity.Category{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? ", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &category)
	if err != nil {
		return nil, 0, err
	}

	return category, total, nil
}

func (conn *categoryConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.Category, string, bool, error) {
	var categories []entity.Category

	query := conn.Db.Model(&entity.Category{}).
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

	if err := query.Find(&categories).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(categories) > limit {
		last := categories[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		categories = categories[:limit]
		hasNext = true
	}

	return categories, nextCursor, hasNext, nil
}

func (conn *categoryConnection) FindWithPaginationCursorProduct(businessId uuid.UUID, pagination request.Pagination) ([]entity.Category, string, bool, error) {
	var categories []entity.Category

	query := conn.Db.Model(&entity.Category{}).
		Where("business_id = ?", businessId)

	query = query.Select(`
    categories.*,
    EXISTS (
        SELECT 1 
        FROM products 
        WHERE products.category_id = categories.id
    	) AS has_product
	`)

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

	if err := query.Find(&categories).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(categories) > limit {
		last := categories[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		categories = categories[:limit]
		hasNext = true
	}

	return categories, nextCursor, hasNext, nil
}
