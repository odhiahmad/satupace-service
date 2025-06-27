package repository

import (
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

// RegistrationRepository interface untuk registrasi bisnis dan user.
type RegistrationRepository interface {
	CreateBusiness(business entity.Business) (entity.Business, error)
	CreateUser(user entity.UserBusiness) error
	CreateMainBranch(branch *entity.BusinessBranch) error
	IsEmailExists(email string) (bool, error)
}

// registrationRepository adalah implementasi dari RegistrationRepository.
type registrationRepository struct {
	db *gorm.DB
}

// NewRegistrationRepository membuat instance baru dari RegistrationRepository.
func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: db}
}

// CreateBusiness menyimpan data bisnis berdasarkan input dari registrasi.
func (r *registrationRepository) CreateBusiness(business entity.Business) (entity.Business, error) {
	if err := r.db.Create(&business).Error; err != nil {
		return business, err
	}
	return business, nil
}

// CreateUser menyimpan data user bisnis ke database.
func (r *registrationRepository) CreateUser(user entity.UserBusiness) error {
	return r.db.Create(&user).Error
}

func (r *registrationRepository) CreateMainBranch(branch *entity.BusinessBranch) error {
	return r.db.Create(branch).Error
}

// IsEmailExists mengecek apakah email sudah digunakan.
func (r *registrationRepository) IsEmailExists(email string) (bool, error) {
	var user entity.UserBusiness
	err := r.db.Where("email = ?", email).First(&user).Error

	switch {
	case err == nil:
		return true, nil
	case err == gorm.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}
