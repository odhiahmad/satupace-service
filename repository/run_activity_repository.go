package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RunActivityRepository interface {
	Create(activity *entity.RunActivity) error
	Update(activity *entity.RunActivity) error
	FindById(id uuid.UUID) (*entity.RunActivity, error)
	FindByUserId(userId uuid.UUID) ([]entity.RunActivity, error)
	FindAll() ([]entity.RunActivity, error)
	Delete(id uuid.UUID) error
	GetUserStats(userId uuid.UUID) (map[string]interface{}, error)
}

type runActivityRepository struct {
	db *gorm.DB
}

func NewRunActivityRepository(db *gorm.DB) RunActivityRepository {
	return &runActivityRepository{db: db}
}

func (r *runActivityRepository) Create(activity *entity.RunActivity) error {
	return r.db.Create(activity).Error
}

func (r *runActivityRepository) Update(activity *entity.RunActivity) error {
	return r.db.Save(activity).Error
}

func (r *runActivityRepository) FindById(id uuid.UUID) (*entity.RunActivity, error) {
	var activity entity.RunActivity
	err := r.db.First(&activity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *runActivityRepository) FindByUserId(userId uuid.UUID) ([]entity.RunActivity, error) {
	var activities []entity.RunActivity
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&activities).Error
	return activities, err
}

func (r *runActivityRepository) FindAll() ([]entity.RunActivity, error) {
	var activities []entity.RunActivity
	err := r.db.Order("created_at DESC").Find(&activities).Error
	return activities, err
}

func (r *runActivityRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.RunActivity{}, "id = ?", id).Error
}

func (r *runActivityRepository) GetUserStats(userId uuid.UUID) (map[string]interface{}, error) {
	var result []map[string]interface{}
	err := r.db.Model(&entity.RunActivity{}).
		Where("user_id = ?", userId).
		Select(
			"COUNT(*) as total_runs",
			"SUM(distance) as total_distance",
			"AVG(avg_pace) as avg_pace",
			"SUM(calories) as total_calories",
		).
		Scan(&result).Error

	if err != nil || len(result) == 0 {
		return nil, err
	}
	return result[0], nil
}
