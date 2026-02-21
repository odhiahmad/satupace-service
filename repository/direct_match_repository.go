package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DirectMatchRepository interface {
	Create(match *entity.DirectMatch) error
	Update(match *entity.DirectMatch) error
	FindById(id uuid.UUID) (*entity.DirectMatch, error)
	FindByUsers(user1Id, user2Id uuid.UUID) (*entity.DirectMatch, error)
	FindUserMatches(userId uuid.UUID) ([]entity.DirectMatch, error)
	FindMatchesByStatus(userId uuid.UUID, status string) ([]entity.DirectMatch, error)
	Delete(id uuid.UUID) error
}

type directMatchRepository struct {
	db *gorm.DB
}

func NewDirectMatchRepository(db *gorm.DB) DirectMatchRepository {
	return &directMatchRepository{db: db}
}

func (r *directMatchRepository) Create(match *entity.DirectMatch) error {
	return r.db.Create(match).Error
}

func (r *directMatchRepository) Update(match *entity.DirectMatch) error {
	return r.db.Save(match).Error
}

func (r *directMatchRepository) FindById(id uuid.UUID) (*entity.DirectMatch, error) {
	var match entity.DirectMatch
	err := r.db.First(&match, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *directMatchRepository) FindByUsers(user1Id, user2Id uuid.UUID) (*entity.DirectMatch, error) {
	var match entity.DirectMatch
	err := r.db.Where(
		"(user_1_id = ? AND user_2_id = ?) OR (user_1_id = ? AND user_2_id = ?)",
		user1Id, user2Id, user2Id, user1Id,
	).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *directMatchRepository) FindUserMatches(userId uuid.UUID) ([]entity.DirectMatch, error) {
	var matches []entity.DirectMatch
	err := r.db.Where("user_1_id = ? OR user_2_id = ?", userId, userId).Order("created_at DESC").Find(&matches).Error
	return matches, err
}

func (r *directMatchRepository) FindMatchesByStatus(userId uuid.UUID, status string) ([]entity.DirectMatch, error) {
	var matches []entity.DirectMatch
	err := r.db.Where(
		"(user_1_id = ? OR user_2_id = ?) AND status = ?",
		userId, userId, status,
	).Order("created_at DESC").Find(&matches).Error
	return matches, err
}

func (r *directMatchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.DirectMatch{}, "id = ?", id).Error
}
