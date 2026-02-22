package repository

import (
	"time"

	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notif *entity.Notification) error
	FindByUserId(userId uuid.UUID, limit, offset int) ([]entity.Notification, error)
	FindUnreadCount(userId uuid.UUID) (int64, error)
	MarkAsRead(ids []uuid.UUID, userId uuid.UUID) error
	MarkAllAsRead(userId uuid.UUID) error
	Delete(id uuid.UUID) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notif *entity.Notification) error {
	return r.db.Create(notif).Error
}

func (r *notificationRepository) FindByUserId(userId uuid.UUID, limit, offset int) ([]entity.Notification, error) {
	var notifs []entity.Notification
	err := r.db.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifs).Error
	return notifs, err
}

func (r *notificationRepository) FindUnreadCount(userId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = false", userId).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(ids []uuid.UUID, userId uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entity.Notification{}).
		Where("id IN ? AND user_id = ?", ids, userId).
		Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}

func (r *notificationRepository) MarkAllAsRead(userId uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = false", userId).
		Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}

func (r *notificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Notification{}, "id = ?", id).Error
}
