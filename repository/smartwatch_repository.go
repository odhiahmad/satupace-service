package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SmartWatchRepository interface {
	// Device management
	CreateDevice(device *entity.SmartWatchDevice) error
	UpdateDevice(device *entity.SmartWatchDevice) error
	FindDeviceById(id uuid.UUID) (*entity.SmartWatchDevice, error)
	FindDevicesByUserId(userId uuid.UUID) ([]entity.SmartWatchDevice, error)
	FindConnectedDevices(userId uuid.UUID) ([]entity.SmartWatchDevice, error)
	DeleteDevice(id uuid.UUID) error

	// Sync management
	CreateSync(sync *entity.SmartWatchSync) error
	CreateSyncBatch(syncs []*entity.SmartWatchSync) error
	FindSyncById(id uuid.UUID) (*entity.SmartWatchSync, error)
	FindSyncsByDeviceId(deviceId uuid.UUID) ([]entity.SmartWatchSync, error)
	FindSyncsByUserId(userId uuid.UUID, limit int) ([]entity.SmartWatchSync, error)
	FindSyncByExternalId(externalId string) (*entity.SmartWatchSync, error)
	GetDeviceStats(deviceId uuid.UUID) (map[string]interface{}, error)
}

type smartWatchRepository struct {
	db *gorm.DB
}

func NewSmartWatchRepository(db *gorm.DB) SmartWatchRepository {
	return &smartWatchRepository{db: db}
}

// ── Device Methods ──

func (r *smartWatchRepository) CreateDevice(device *entity.SmartWatchDevice) error {
	return r.db.Create(device).Error
}

func (r *smartWatchRepository) UpdateDevice(device *entity.SmartWatchDevice) error {
	return r.db.Save(device).Error
}

func (r *smartWatchRepository) FindDeviceById(id uuid.UUID) (*entity.SmartWatchDevice, error) {
	var device entity.SmartWatchDevice
	err := r.db.First(&device, "id = ?", id).Error
	return &device, err
}

func (r *smartWatchRepository) FindDevicesByUserId(userId uuid.UUID) ([]entity.SmartWatchDevice, error) {
	var devices []entity.SmartWatchDevice
	err := r.db.Where("user_id = ?", userId).Order("connected_at DESC").Find(&devices).Error
	return devices, err
}

func (r *smartWatchRepository) FindConnectedDevices(userId uuid.UUID) ([]entity.SmartWatchDevice, error) {
	var devices []entity.SmartWatchDevice
	err := r.db.Where("user_id = ? AND is_connected = ?", userId, true).Find(&devices).Error
	return devices, err
}

func (r *smartWatchRepository) DeleteDevice(id uuid.UUID) error {
	return r.db.Delete(&entity.SmartWatchDevice{}, "id = ?", id).Error
}

// ── Sync Methods ──

func (r *smartWatchRepository) CreateSync(sync *entity.SmartWatchSync) error {
	return r.db.Create(sync).Error
}

func (r *smartWatchRepository) CreateSyncBatch(syncs []*entity.SmartWatchSync) error {
	return r.db.CreateInBatches(syncs, 50).Error
}

func (r *smartWatchRepository) FindSyncById(id uuid.UUID) (*entity.SmartWatchSync, error) {
	var sync entity.SmartWatchSync
	err := r.db.First(&sync, "id = ?", id).Error
	return &sync, err
}

func (r *smartWatchRepository) FindSyncsByDeviceId(deviceId uuid.UUID) ([]entity.SmartWatchSync, error) {
	var syncs []entity.SmartWatchSync
	err := r.db.Where("device_id = ?", deviceId).Order("activity_date DESC").Find(&syncs).Error
	return syncs, err
}

func (r *smartWatchRepository) FindSyncsByUserId(userId uuid.UUID, limit int) ([]entity.SmartWatchSync, error) {
	var syncs []entity.SmartWatchSync
	q := r.db.Where("user_id = ?", userId).Order("activity_date DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&syncs).Error
	return syncs, err
}

func (r *smartWatchRepository) FindSyncByExternalId(externalId string) (*entity.SmartWatchSync, error) {
	var sync entity.SmartWatchSync
	err := r.db.First(&sync, "external_id = ?", externalId).Error
	return &sync, err
}

func (r *smartWatchRepository) GetDeviceStats(deviceId uuid.UUID) (map[string]interface{}, error) {
	var result []map[string]interface{}
	err := r.db.Model(&entity.SmartWatchSync{}).
		Where("device_id = ? AND status = ?", deviceId, "synced").
		Select(
			"COUNT(*) as total_activities",
			"COALESCE(SUM(distance), 0) as total_distance",
			"COALESCE(SUM(duration), 0) as total_duration",
			"COALESCE(AVG(avg_pace), 0) as avg_pace",
			"COALESCE(AVG(NULLIF(avg_heart_rate, 0)), 0) as avg_heart_rate",
			"COALESCE(SUM(calories), 0) as total_calories",
			"COALESCE(SUM(elevation_gain), 0) as total_elevation",
			"MAX(activity_date) as last_activity_date",
		).
		Scan(&result).Error

	if err != nil || len(result) == 0 {
		return nil, err
	}
	return result[0], nil
}
