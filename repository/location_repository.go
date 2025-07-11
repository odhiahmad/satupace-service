package repository

import (
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type LocationRepository interface {
	GetProvinces() ([]entity.Province, error)
	GetCitiesByProvinceID(provinceID int) ([]entity.City, error)
	GetDistrictsByCityID(cityID int) ([]entity.District, error)
	GetVillagesByDistrictID(districtID int) ([]entity.Village, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db}
}

func (r *locationRepository) GetProvinces() ([]entity.Province, error) {
	var provinces []entity.Province
	err := r.db.Find(&provinces).Error
	return provinces, err
}

func (r *locationRepository) GetCitiesByProvinceID(provinceID int) ([]entity.City, error) {
	var cities []entity.City
	err := r.db.Where("province_id = ?", provinceID).Find(&cities).Error
	return cities, err
}

func (r *locationRepository) GetDistrictsByCityID(cityID int) ([]entity.District, error) {
	var districts []entity.District
	err := r.db.Where("city_id = ?", cityID).Find(&districts).Error
	return districts, err
}

func (r *locationRepository) GetVillagesByDistrictID(districtID int) ([]entity.Village, error) {
	var village []entity.Village
	err := r.db.Where("district_id = ?", districtID).Find(&village).Error
	return village, err
}
