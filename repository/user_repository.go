package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	Update(user *entity.User) error
	FindById(id uuid.UUID) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByPhone(phone string) (*entity.User, error)
	FindAll() ([]entity.User, error)
	Delete(id uuid.UUID) error
	FindByEmailOrPhone(identifier string) (*entity.User, error)
	IsDuplicateEmail(email string) bool
	IsDuplicatePhone(phone string) bool
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindById(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhone(phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "phone_number = ?", phone).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}

func (r *userRepository) FindByEmailOrPhone(identifier string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ? OR phone_number = ?", identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) IsDuplicateEmail(email string) bool {
	var count int64
	r.db.Model(&entity.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *userRepository) IsDuplicatePhone(phone string) bool {
	var count int64
	r.db.Model(&entity.User{}).Where("phone_number = ?", phone).Count(&count)
	return count > 0
}
