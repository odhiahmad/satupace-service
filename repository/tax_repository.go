package repository

import (
	"errors"

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
}

type taxConnection struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) TaxRepository {
	return &taxConnection{db}
}

func (conn *taxConnection) Create(tax entity.Tax) (entity.Tax, error) {
	// Buat data tax
	err := conn.db.Create(&tax).Error
	if err != nil {
		return entity.Tax{}, err
	}

	// Ambil ulang dengan preload Business
	err = conn.db.Preload("Business").First(&tax, tax.Id).Error
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

	err = conn.db.Preload("Business").First(&tax, tax.Id).Error

	return tax, err
}

func (conn *taxConnection) Delete(tax entity.Tax) error {
	return conn.db.Preload("Business").Delete(&tax).Error
}

func (conn *taxConnection) FindById(taxId int) (taxes entity.Tax, err error) {
	var tax entity.Tax
	result := conn.db.Preload("Business").Find(&tax, taxId)
	if result != nil {
		return tax, nil
	} else {
		return tax, errors.New("tag is not found")
	}
}

func (conn *taxConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Tax, int64, error) {
	var tax []entity.Tax
	var total int64

	// Base query dengan preload relasi
	baseQuery := conn.db.Model(&entity.Tax{}).
		Preload("Business").
		Where("business_id = ?", businessId)

	// Search filter
	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	// Hitung total sebelum paginasi
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &tax)
	if err != nil {
		return nil, 0, err
	}

	return tax, total, nil
}
