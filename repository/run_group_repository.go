package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RunGroupRepository interface {
	Create(group *entity.RunGroup) error
	Update(group *entity.RunGroup) error
	FindById(id uuid.UUID) (*entity.RunGroup, error)
	FindAll() ([]entity.RunGroup, error)
	FindByStatus(status string) ([]entity.RunGroup, error)
	Delete(id uuid.UUID) error
	FindByCreatedBy(userId uuid.UUID) ([]entity.RunGroup, error)
	GetMemberCount(groupId uuid.UUID) (int64, error)
}

type runGroupRepository struct {
	db *gorm.DB
}

func NewRunGroupRepository(db *gorm.DB) RunGroupRepository {
	return &runGroupRepository{db: db}
}

func (r *runGroupRepository) Create(group *entity.RunGroup) error {
	return r.db.Create(group).Error
}

func (r *runGroupRepository) Update(group *entity.RunGroup) error {
	return r.db.Save(group).Error
}

func (r *runGroupRepository) FindById(id uuid.UUID) (*entity.RunGroup, error) {
	var group entity.RunGroup
	err := r.db.First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *runGroupRepository) FindAll() ([]entity.RunGroup, error) {
	var groups []entity.RunGroup
	err := r.db.Find(&groups).Error
	return groups, err
}

func (r *runGroupRepository) FindByStatus(status string) ([]entity.RunGroup, error) {
	var groups []entity.RunGroup
	err := r.db.Where("status = ?", status).Find(&groups).Error
	return groups, err
}

func (r *runGroupRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.RunGroup{}, "id = ?", id).Error
}

func (r *runGroupRepository) FindByCreatedBy(userId uuid.UUID) ([]entity.RunGroup, error) {
	var groups []entity.RunGroup
	err := r.db.Where("created_by = ?", userId).Find(&groups).Error
	return groups, err
}

func (r *runGroupRepository) GetMemberCount(groupId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.RunGroupMember{}).Where("group_id = ? AND status = ?", groupId, "joined").Count(&count).Error
	return count, err
}
