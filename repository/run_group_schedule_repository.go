package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RunGroupScheduleRepository interface {
	Create(schedule *entity.RunGroupSchedule) error
	Update(schedule *entity.RunGroupSchedule) error
	FindById(id uuid.UUID) (*entity.RunGroupSchedule, error)
	FindByGroupId(groupId uuid.UUID) ([]entity.RunGroupSchedule, error)
	CountByGroupId(groupId uuid.UUID) (int64, error)
	Delete(id uuid.UUID) error
}

type runGroupScheduleRepository struct {
	db *gorm.DB
}

func NewRunGroupScheduleRepository(db *gorm.DB) RunGroupScheduleRepository {
	return &runGroupScheduleRepository{db: db}
}

func (r *runGroupScheduleRepository) Create(schedule *entity.RunGroupSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *runGroupScheduleRepository) Update(schedule *entity.RunGroupSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *runGroupScheduleRepository) FindById(id uuid.UUID) (*entity.RunGroupSchedule, error) {
	var schedule entity.RunGroupSchedule
	err := r.db.First(&schedule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *runGroupScheduleRepository) FindByGroupId(groupId uuid.UUID) ([]entity.RunGroupSchedule, error) {
	var schedules []entity.RunGroupSchedule
	err := r.db.Where("group_id = ?", groupId).
		Order("day_of_week ASC, start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *runGroupScheduleRepository) CountByGroupId(groupId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.RunGroupSchedule{}).Where("group_id = ? AND is_active = ?", groupId, true).Count(&count).Error
	return count, err
}

func (r *runGroupScheduleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.RunGroupSchedule{}, "id = ?", id).Error
}
