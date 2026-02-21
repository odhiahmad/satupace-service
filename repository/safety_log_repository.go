package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SafetyLogRepository interface {
	Create(log *entity.SafetyLog) error
	FindById(id uuid.UUID) (*entity.SafetyLog, error)
	FindByUserId(userId uuid.UUID) ([]entity.SafetyLog, error)
	FindByMatchId(matchId uuid.UUID) ([]entity.SafetyLog, error)
	FindByStatus(status string) ([]entity.SafetyLog, error)
	Delete(id uuid.UUID) error
	GetUserSafetyLogs(userId uuid.UUID) ([]entity.SafetyLog, error)
	CountReportsByTarget(targetUserId uuid.UUID) (int64, error)
	DB() *gorm.DB
}

type safetyLogRepository struct {
	db *gorm.DB
}

func NewSafetyLogRepository(db *gorm.DB) SafetyLogRepository {
	return &safetyLogRepository{db: db}
}

func (r *safetyLogRepository) Create(log *entity.SafetyLog) error {
	return r.db.Create(log).Error
}

func (r *safetyLogRepository) FindById(id uuid.UUID) (*entity.SafetyLog, error) {
	var log entity.SafetyLog
	err := r.db.First(&log, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *safetyLogRepository) FindByUserId(userId uuid.UUID) ([]entity.SafetyLog, error) {
	var logs []entity.SafetyLog
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *safetyLogRepository) FindByMatchId(matchId uuid.UUID) ([]entity.SafetyLog, error) {
	var logs []entity.SafetyLog
	err := r.db.Where("match_id = ?", matchId).Find(&logs).Error
	return logs, err
}

func (r *safetyLogRepository) FindByStatus(status string) ([]entity.SafetyLog, error) {
	var logs []entity.SafetyLog
	err := r.db.Where("status = ?", status).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *safetyLogRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.SafetyLog{}, "id = ?", id).Error
}

func (r *safetyLogRepository) GetUserSafetyLogs(userId uuid.UUID) ([]entity.SafetyLog, error) {
	var logs []entity.SafetyLog
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *safetyLogRepository) CountReportsByTarget(targetUserId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.SafetyLog{}).Where("match_id = ? AND status = ?", targetUserId, "reported").Count(&count).Error
	return count, err
}

func (r *safetyLogRepository) DB() *gorm.DB {
	return r.db
}
