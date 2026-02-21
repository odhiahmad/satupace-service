package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BiometricRepository interface {
	Create(biometric *entity.UserBiometric) error
	Update(biometric *entity.UserBiometric) error
	FindByCredentialId(credentialId string) (*entity.UserBiometric, error)
	FindByUserId(userId uuid.UUID) ([]entity.UserBiometric, error)
	Delete(id uuid.UUID) error
}

type biometricRepository struct {
	db *gorm.DB
}

func NewBiometricRepository(db *gorm.DB) BiometricRepository {
	return &biometricRepository{db: db}
}

func (r *biometricRepository) Create(biometric *entity.UserBiometric) error {
	return r.db.Create(biometric).Error
}

func (r *biometricRepository) Update(biometric *entity.UserBiometric) error {
	return r.db.Save(biometric).Error
}

func (r *biometricRepository) FindByCredentialId(credentialId string) (*entity.UserBiometric, error) {
	var biometric entity.UserBiometric
	err := r.db.Where("credential_id = ? AND is_active = true", credentialId).First(&biometric).Error
	if err != nil {
		return nil, err
	}
	return &biometric, nil
}

func (r *biometricRepository) FindByUserId(userId uuid.UUID) ([]entity.UserBiometric, error) {
	var biometrics []entity.UserBiometric
	err := r.db.Where("user_id = ? AND is_active = true", userId).Find(&biometrics).Error
	return biometrics, err
}

func (r *biometricRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.UserBiometric{}, "id = ?", id).Error
}
