package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPhotoRepository interface {
	Create(photo *entity.UserPhoto) error
	Update(photo *entity.UserPhoto) error
	FindById(id uuid.UUID) (*entity.UserPhoto, error)
	FindByUserId(userId uuid.UUID) ([]entity.UserPhoto, error)
	FindPrimaryPhoto(userId uuid.UUID) (*entity.UserPhoto, error)
	FindVerificationPhoto(userId uuid.UUID) (*entity.UserPhoto, error)
	Delete(id uuid.UUID) error
	DeleteByUserId(userId uuid.UUID) error
}

type userPhotoRepository struct {
	db *gorm.DB
}

func NewUserPhotoRepository(db *gorm.DB) UserPhotoRepository {
	return &userPhotoRepository{db: db}
}

func (r *userPhotoRepository) Create(photo *entity.UserPhoto) error {
	return r.db.Create(photo).Error
}

func (r *userPhotoRepository) Update(photo *entity.UserPhoto) error {
	return r.db.Save(photo).Error
}

func (r *userPhotoRepository) FindById(id uuid.UUID) (*entity.UserPhoto, error) {
	var photo entity.UserPhoto
	err := r.db.First(&photo, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *userPhotoRepository) FindByUserId(userId uuid.UUID) ([]entity.UserPhoto, error) {
	var photos []entity.UserPhoto
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&photos).Error
	return photos, err
}

func (r *userPhotoRepository) FindPrimaryPhoto(userId uuid.UUID) (*entity.UserPhoto, error) {
	var photo entity.UserPhoto
	err := r.db.Where("user_id = ? AND is_primary = ?", userId, true).First(&photo).Error
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *userPhotoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.UserPhoto{}, "id = ?", id).Error
}

func (r *userPhotoRepository) FindVerificationPhoto(userId uuid.UUID) (*entity.UserPhoto, error) {
	var photo entity.UserPhoto
	err := r.db.Where("user_id = ? AND type = 'verification'", userId).Order("created_at DESC").First(&photo).Error
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *userPhotoRepository) DeleteByUserId(userId uuid.UUID) error {
	return r.db.Delete(&entity.UserPhoto{}, "user_id = ?", userId).Error
}
