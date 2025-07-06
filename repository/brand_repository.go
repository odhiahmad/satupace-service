package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BrandRepository interface {
	Create(brand entity.Brand) (entity.Brand, error)
	Update(brand entity.Brand) (entity.Brand, error)
	Delete(brand entity.Brand) error
	FindById(brandId int) (brandes entity.Brand, err error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Brand, int64, error)
}

type brandConnection struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandConnection{db}
}

func (conn *brandConnection) Create(brand entity.Brand) (entity.Brand, error) {
	// Buat data brand
	err := conn.db.Create(&brand).Error
	if err != nil {
		return entity.Brand{}, err
	}

	// Ambil ulang dengan preload Business
	err = conn.db.First(&brand, brand.Id).Error
	if err != nil {
		return entity.Brand{}, err
	}

	return brand, nil
}

func (conn *brandConnection) Update(brand entity.Brand) (entity.Brand, error) {
	err := conn.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&brand).Error
	if err != nil {
		return entity.Brand{}, err
	}

	err = conn.db.First(&brand, brand.Id).Error

	return brand, err
}

func (conn *brandConnection) Delete(brand entity.Brand) error {
	return conn.db.Delete(&brand).Error
}

func (conn *brandConnection) FindById(brandId int) (brandes entity.Brand, err error) {
	var brand entity.Brand
	result := conn.db.Find(&brand, brandId)
	if result != nil {
		return brand, nil
	} else {
		return brand, errors.New("tag is not found")
	}
}

func (conn *brandConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Brand, int64, error) {
	var brand []entity.Brand
	var total int64

	// Base query dengan preload relasi
	baseQuery := conn.db.Model(&entity.Brand{}).
		Where("business_id = ?", businessId)

	// Search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	// Hitung total sebelum paginasi
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &brand)
	if err != nil {
		return nil, 0, err
	}

	return brand, total, nil
}
