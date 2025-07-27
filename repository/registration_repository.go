package repository

import (
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

// RegistrationRepository interface untuk registrasi bisnis dan user.
type RegistrationRepository interface {
	CreateBusiness(business entity.Business) (entity.Business, error)
	CreateUser(user entity.UserBusiness) (entity.UserBusiness, error)
	CreateDefaultCategory(category entity.Category) error
	IsEmailExists(email string) (bool, error)
	IsPhoneNumberExists(phone string) (bool, error)
}

type registrationConnection struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationConnection{db: db}
}

func (conn *registrationConnection) CreateBusiness(business entity.Business) (entity.Business, error) {
	if err := conn.db.Create(&business).Error; err != nil {
		return business, err
	}
	return business, nil
}

func (conn *registrationConnection) CreateUser(user entity.UserBusiness) (entity.UserBusiness, error) {
	if err := conn.db.Create(&user).Error; err != nil {
		return entity.UserBusiness{}, err
	}

	var savedUser entity.UserBusiness
	if err := conn.db.
		First(&savedUser, user.Id).Error; err != nil {
		return entity.UserBusiness{}, err
	}

	return savedUser, nil
}

func (conn *registrationConnection) CreateDefaultCategory(category entity.Category) error {
	return conn.db.Create(&category).Error
}

func (conn *registrationConnection) IsEmailExists(email string) (bool, error) {
	var user entity.UserBusiness
	err := conn.db.Where("email = ?", email).First(&user).Error

	switch {
	case err == nil:
		return true, nil
	case err == gorm.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (conn *registrationConnection) IsPhoneNumberExists(phone string) (bool, error) {
	var user entity.UserBusiness
	err := conn.db.Where("phone_number = ?", phone).First(&user).Error

	switch {
	case err == nil:
		return true, nil
	case err == gorm.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}
