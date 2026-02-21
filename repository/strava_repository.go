package repository

import (
	"run-sync/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StravaRepository interface {
	// Connection
	CreateConnection(conn *entity.StravaConnection) error
	UpdateConnection(conn *entity.StravaConnection) error
	FindConnectionByUserId(userId uuid.UUID) (*entity.StravaConnection, error)
	FindConnectionByAthleteId(athleteId int64) (*entity.StravaConnection, error)
	DeleteConnection(id uuid.UUID) error

	// Activity
	CreateActivity(activity *entity.StravaActivity) error
	CreateActivitiesBatch(activities []*entity.StravaActivity) error
	FindActivitiesByUserId(userId uuid.UUID, limit int) ([]entity.StravaActivity, error)
	FindActivityByStravaId(stravaId int64) (*entity.StravaActivity, error)
	GetSyncStats(userId uuid.UUID) (map[string]interface{}, error)
}

type stravaRepository struct {
	db *gorm.DB
}

func NewStravaRepository(db *gorm.DB) StravaRepository {
	return &stravaRepository{db: db}
}

// ── Connection ──

func (r *stravaRepository) CreateConnection(conn *entity.StravaConnection) error {
	return r.db.Create(conn).Error
}

func (r *stravaRepository) UpdateConnection(conn *entity.StravaConnection) error {
	return r.db.Save(conn).Error
}

func (r *stravaRepository) FindConnectionByUserId(userId uuid.UUID) (*entity.StravaConnection, error) {
	var conn entity.StravaConnection
	err := r.db.First(&conn, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *stravaRepository) FindConnectionByAthleteId(athleteId int64) (*entity.StravaConnection, error) {
	var conn entity.StravaConnection
	err := r.db.First(&conn, "athlete_id = ?", athleteId).Error
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *stravaRepository) DeleteConnection(id uuid.UUID) error {
	return r.db.Delete(&entity.StravaConnection{}, "id = ?", id).Error
}

// ── Activity ──

func (r *stravaRepository) CreateActivity(activity *entity.StravaActivity) error {
	return r.db.Create(activity).Error
}

func (r *stravaRepository) CreateActivitiesBatch(activities []*entity.StravaActivity) error {
	return r.db.CreateInBatches(activities, 50).Error
}

func (r *stravaRepository) FindActivitiesByUserId(userId uuid.UUID, limit int) ([]entity.StravaActivity, error) {
	var activities []entity.StravaActivity
	q := r.db.Where("user_id = ?", userId).Order("start_date DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&activities).Error
	return activities, err
}

func (r *stravaRepository) FindActivityByStravaId(stravaId int64) (*entity.StravaActivity, error) {
	var activity entity.StravaActivity
	err := r.db.First(&activity, "strava_id = ?", stravaId).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *stravaRepository) GetSyncStats(userId uuid.UUID) (map[string]interface{}, error) {
	var result []map[string]interface{}
	err := r.db.Model(&entity.StravaActivity{}).
		Where("user_id = ? AND status = ?", userId, "synced").
		Select(
			"COUNT(*) as total_activities",
			"COALESCE(SUM(distance), 0) as total_distance",
			"COALESCE(SUM(moving_time), 0) as total_duration",
			"COALESCE(AVG(CASE WHEN average_speed > 0 THEN average_speed END), 0) as avg_speed",
			"COALESCE(AVG(NULLIF(average_heartrate, 0)), 0) as avg_heartrate",
			"COALESCE(SUM(calories), 0) as total_calories",
			"COALESCE(SUM(total_elevation), 0) as total_elevation",
			"MAX(start_date) as last_activity_date",
		).
		Scan(&result).Error

	if err != nil || len(result) == 0 {
		return nil, err
	}
	return result[0], nil
}
