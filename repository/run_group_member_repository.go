package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RunGroupMemberRepository interface {
	Create(member *entity.RunGroupMember) error
	Update(member *entity.RunGroupMember) error
	FindById(id uuid.UUID) (*entity.RunGroupMember, error)
	FindByGroupAndUser(groupId, userId uuid.UUID) (*entity.RunGroupMember, error)
	FindByGroupId(groupId uuid.UUID) ([]entity.RunGroupMember, error)
	FindByUserId(userId uuid.UUID) ([]entity.RunGroupMember, error)
	Delete(id uuid.UUID) error
	GetMembers(groupId uuid.UUID, status string) ([]entity.RunGroupMember, error)
	DB() *gorm.DB
}

type runGroupMemberRepository struct {
	db *gorm.DB
}

func NewRunGroupMemberRepository(db *gorm.DB) RunGroupMemberRepository {
	return &runGroupMemberRepository{db: db}
}

func (r *runGroupMemberRepository) Create(member *entity.RunGroupMember) error {
	return r.db.Create(member).Error
}

func (r *runGroupMemberRepository) Update(member *entity.RunGroupMember) error {
	return r.db.Save(member).Error
}

func (r *runGroupMemberRepository) FindById(id uuid.UUID) (*entity.RunGroupMember, error) {
	var member entity.RunGroupMember
	err := r.db.First(&member, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *runGroupMemberRepository) FindByGroupAndUser(groupId, userId uuid.UUID) (*entity.RunGroupMember, error) {
	var member entity.RunGroupMember
	err := r.db.First(&member, "group_id = ? AND user_id = ?", groupId, userId).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *runGroupMemberRepository) FindByGroupId(groupId uuid.UUID) ([]entity.RunGroupMember, error) {
	var members []entity.RunGroupMember
	err := r.db.Where("group_id = ?", groupId).Find(&members).Error
	return members, err
}

func (r *runGroupMemberRepository) FindByUserId(userId uuid.UUID) ([]entity.RunGroupMember, error) {
	var members []entity.RunGroupMember
	err := r.db.Where("user_id = ?", userId).Find(&members).Error
	return members, err
}

func (r *runGroupMemberRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.RunGroupMember{}, "id = ?", id).Error
}

func (r *runGroupMemberRepository) GetMembers(groupId uuid.UUID, status string) ([]entity.RunGroupMember, error) {
	var members []entity.RunGroupMember
	err := r.db.Where("group_id = ? AND status = ?", groupId, status).Find(&members).Error
	return members, err
}

func (r *runGroupMemberRepository) DB() *gorm.DB {
	return r.db
}
