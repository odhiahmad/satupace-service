package repository

import (
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type RegistrationRepository interface {
	CreateBusiness(business entity.Business) (entity.Business, error)
	CreateUser(user entity.UserBusiness) error
	IsEmailExists(email string) (bool, error)
}

type registrationConnection struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationConnection{db: db}
}

func (r *registrationConnection) CreateBusiness(business entity.Business) (entity.Business, error) {
	if err := r.db.Create(&business).Error; err != nil {
		return business, err
	}
	return business, nil
}

func (r *registrationConnection) CreateUser(user entity.UserBusiness) error {
	return r.db.Create(&user).Error
}

func (r *registrationConnection) IsEmailExists(email string) (bool, error) {
	var user entity.UserBusiness
	err := r.db.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
