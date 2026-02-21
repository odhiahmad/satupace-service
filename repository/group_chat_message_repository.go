package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupChatMessageRepository interface {
	Create(message *entity.GroupChatMessage) error
	FindById(id uuid.UUID) (*entity.GroupChatMessage, error)
	FindByGroupId(groupId uuid.UUID) ([]entity.GroupChatMessage, error)
	FindBySenderId(userId uuid.UUID) ([]entity.GroupChatMessage, error)
	Delete(id uuid.UUID) error
	DeleteByGroupId(groupId uuid.UUID) error
}

type groupChatMessageRepository struct {
	db *gorm.DB
}

func NewGroupChatMessageRepository(db *gorm.DB) GroupChatMessageRepository {
	return &groupChatMessageRepository{db: db}
}

func (r *groupChatMessageRepository) Create(message *entity.GroupChatMessage) error {
	return r.db.Create(message).Error
}

func (r *groupChatMessageRepository) FindById(id uuid.UUID) (*entity.GroupChatMessage, error) {
	var message entity.GroupChatMessage
	err := r.db.First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *groupChatMessageRepository) FindByGroupId(groupId uuid.UUID) ([]entity.GroupChatMessage, error) {
	var messages []entity.GroupChatMessage
	err := r.db.Where("group_id = ?", groupId).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *groupChatMessageRepository) FindBySenderId(userId uuid.UUID) ([]entity.GroupChatMessage, error) {
	var messages []entity.GroupChatMessage
	err := r.db.Where("sender_id = ?", userId).Order("created_at DESC").Find(&messages).Error
	return messages, err
}

func (r *groupChatMessageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.GroupChatMessage{}, "id = ?", id).Error
}

func (r *groupChatMessageRepository) DeleteByGroupId(groupId uuid.UUID) error {
	return r.db.Delete(&entity.GroupChatMessage{}, "group_id = ?", groupId).Error
}
