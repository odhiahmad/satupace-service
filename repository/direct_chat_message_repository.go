package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DirectChatMessageRepository interface {
	Create(message *entity.DirectChatMessage) error
	FindById(id uuid.UUID) (*entity.DirectChatMessage, error)
	FindByMatchId(matchId uuid.UUID) ([]entity.DirectChatMessage, error)
	FindBySenderId(userId uuid.UUID) ([]entity.DirectChatMessage, error)
	Delete(id uuid.UUID) error
	DeleteByMatchId(matchId uuid.UUID) error
}

type directChatMessageRepository struct {
	db *gorm.DB
}

func NewDirectChatMessageRepository(db *gorm.DB) DirectChatMessageRepository {
	return &directChatMessageRepository{db: db}
}

func (r *directChatMessageRepository) Create(message *entity.DirectChatMessage) error {
	return r.db.Create(message).Error
}

func (r *directChatMessageRepository) FindById(id uuid.UUID) (*entity.DirectChatMessage, error) {
	var message entity.DirectChatMessage
	err := r.db.First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *directChatMessageRepository) FindByMatchId(matchId uuid.UUID) ([]entity.DirectChatMessage, error) {
	var messages []entity.DirectChatMessage
	err := r.db.Where("match_id = ?", matchId).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *directChatMessageRepository) FindBySenderId(userId uuid.UUID) ([]entity.DirectChatMessage, error) {
	var messages []entity.DirectChatMessage
	err := r.db.Where("sender_id = ?", userId).Order("created_at DESC").Find(&messages).Error
	return messages, err
}

func (r *directChatMessageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.DirectChatMessage{}, "id = ?", id).Error
}

func (r *directChatMessageRepository) DeleteByMatchId(matchId uuid.UUID) error {
	return r.db.Delete(&entity.DirectChatMessage{}, "match_id = ?", matchId).Error
}
