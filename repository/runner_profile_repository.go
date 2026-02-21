package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RunnerProfileRepository interface {
	Create(profile *entity.RunnerProfile) error
	Update(profile *entity.RunnerProfile) error
	FindById(id uuid.UUID) (*entity.RunnerProfile, error)
	FindByUserId(userId uuid.UUID) (*entity.RunnerProfile, error)
	FindAll() ([]entity.RunnerProfile, error)
	Delete(id uuid.UUID) error
}

type runnerProfileRepository struct {
	db *gorm.DB
}

func NewRunnerProfileRepository(db *gorm.DB) RunnerProfileRepository {
	return &runnerProfileRepository{db: db}
}

func (r *runnerProfileRepository) Create(profile *entity.RunnerProfile) error {
	return r.db.Create(profile).Error
}

func (r *runnerProfileRepository) Update(profile *entity.RunnerProfile) error {
	return r.db.Save(profile).Error
}

func (r *runnerProfileRepository) FindById(id uuid.UUID) (*entity.RunnerProfile, error) {
	var profile entity.RunnerProfile
	err := r.db.Preload("User").First(&profile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *runnerProfileRepository) FindByUserId(userId uuid.UUID) (*entity.RunnerProfile, error) {
	var profile entity.RunnerProfile
	err := r.db.Preload("User").First(&profile, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *runnerProfileRepository) FindAll() ([]entity.RunnerProfile, error) {
	var profiles []entity.RunnerProfile
	err := r.db.Preload("User").Find(&profiles).Error
	return profiles, err
}

func (r *runnerProfileRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.RunnerProfile{}, "id = ?", id).Error
}
