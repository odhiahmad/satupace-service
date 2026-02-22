package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDeviceTokenRepository interface {
	Upsert(token *entity.UserDeviceToken) error
	FindByUserId(userId uuid.UUID) ([]entity.UserDeviceToken, error)
	DeleteByToken(fcmToken string) error
}

type userDeviceTokenRepository struct {
	db *gorm.DB
}

func NewUserDeviceTokenRepository(db *gorm.DB) UserDeviceTokenRepository {
	return &userDeviceTokenRepository{db: db}
}

// Upsert: kalau token sudah ada, update userId dan platform-nya.
// Berguna saat user logout lalu login di device yang sama.
func (r *userDeviceTokenRepository) Upsert(token *entity.UserDeviceToken) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "fcm_token"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "platform", "is_active", "updated_at"}),
	}).Create(token).Error
}

func (r *userDeviceTokenRepository) FindByUserId(userId uuid.UUID) ([]entity.UserDeviceToken, error) {
	var tokens []entity.UserDeviceToken
	err := r.db.Where("user_id = ? AND is_active = true", userId).Find(&tokens).Error
	return tokens, err
}

// DeleteByToken: dipanggil saat user logout dari device tersebut.
func (r *userDeviceTokenRepository) DeleteByToken(fcmToken string) error {
	return r.db.Where("fcm_token = ?", fcmToken).Delete(&entity.UserDeviceToken{}).Error
}
