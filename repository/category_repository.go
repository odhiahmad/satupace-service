package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	InsertCategory(category entity.Category) (entity.Category, error)
	UpdateCategory(category entity.Category) (entity.Category, error)
	FindById(categoryId int) (entity.Category, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Category, int64, error)
	Delete(categoryId int) error
}

type categoryConnection struct {
	Db *gorm.DB
}

func NewCategoryRepository(Db *gorm.DB) CategoryRepository {
	return &categoryConnection{Db: Db}
}

func (conn *categoryConnection) InsertCategory(category entity.Category) (entity.Category, error) {
	// Validasi keberadaan business
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

	// Ambil ulang dengan preload relasi Business
	err = conn.Db.Preload("Business").First(&category, category.Id).Error
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

	// Ambil ulang dengan preload relasi Business
	err = conn.Db.Preload("Business").First(&category, category.Id).Error
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (conn *categoryConnection) FindById(categoryId int) (entity.Category, error) {
	var category entity.Category
	result := conn.Db.First(&category, categoryId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return category, errors.New("category not found")
	}
	return category, result.Error
}

func (conn *categoryConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Category, int64, error) {
	var category []entity.Category
	var total int64

	// Base query dengan preload relasi
	baseQuery := conn.Db.Model(&entity.Category{}).
		Preload("Business").
		Where("business_id = ?", businessId)

	// Search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE", search)
	}

	// Hitung total sebelum paginasi
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &category)
	if err != nil {
		return nil, 0, err
	}

	return category, total, nil
}

func (conn *categoryConnection) Delete(categoryId int) error {
	result := conn.Db.Delete(&entity.Category{}, categoryId)
	return result.Error
}
