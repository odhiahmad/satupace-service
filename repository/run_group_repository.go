package repository

import (
	"math"
	"run-sync/data/request"
	"run-sync/entity"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RunGroupRepository interface {
	Create(group *entity.RunGroup) error
	Update(group *entity.RunGroup) error
	FindById(id uuid.UUID) (*entity.RunGroup, error)
	FindAll(filter request.RunGroupFilterRequest) ([]entity.RunGroup, error)
	FindByStatus(status string) ([]entity.RunGroup, error)
	Delete(id uuid.UUID) error
	FindByCreatedBy(userId uuid.UUID) ([]entity.RunGroup, error)
	FindByMembership(userId uuid.UUID) ([]entity.RunGroup, []string, error)
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
	err := r.db.Preload("Schedules").First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *runGroupRepository) FindAll(filter request.RunGroupFilterRequest) ([]entity.RunGroup, error) {
	var groups []entity.RunGroup
	query := r.db.Model(&entity.RunGroup{})

	// Filter by status
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Filter by women only
	if filter.WomenOnly == "true" {
		query = query.Where("is_women_only = ?", true)
	}

	// Filter by pace range (overlap: group's pace range intersects with filter's pace range)
	if filter.MinPace != "" {
		if minPace, err := strconv.ParseFloat(filter.MinPace, 64); err == nil {
			query = query.Where("max_pace >= ?", minPace)
		}
	}
	if filter.MaxPace != "" {
		if maxPace, err := strconv.ParseFloat(filter.MaxPace, 64); err == nil {
			query = query.Where("min_pace <= ?", maxPace)
		}
	}

	// Filter by max preferred distance
	if filter.MaxDistance != "" {
		if maxDist, err := strconv.Atoi(filter.MaxDistance); err == nil {
			query = query.Where("preferred_distance <= ?", maxDist)
		}
	}

	// Filter by location radius (Haversine)
	if filter.Latitude != "" && filter.Longitude != "" && filter.RadiusKm != "" {
		lat, errLat := strconv.ParseFloat(filter.Latitude, 64)
		lng, errLng := strconv.ParseFloat(filter.Longitude, 64)
		radiusKm, errR := strconv.ParseFloat(filter.RadiusKm, 64)
		if errLat == nil && errLng == nil && errR == nil {
			// Approximate bounding box for performance
			latDelta := radiusKm / 111.0
			lngDelta := radiusKm / (111.0 * math.Cos(lat*math.Pi/180.0))

			query = query.Where("latitude BETWEEN ? AND ?", lat-latDelta, lat+latDelta)
			query = query.Where("longitude BETWEEN ? AND ?", lng-lngDelta, lng+lngDelta)

			// Precise Haversine filter
			query = query.Where(
				"( 6371 * acos( cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)) ) ) <= ?",
				lat, lng, lat, radiusKm,
			)
		}
	}

	err := query.Preload("Schedules").Order("created_at DESC").Find(&groups).Error
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
	err := r.db.Preload("Schedules").Where("created_by = ?", userId).Find(&groups).Error
	return groups, err
}

func (r *runGroupRepository) FindByMembership(userId uuid.UUID) ([]entity.RunGroup, []string, error) {
	type result struct {
		entity.RunGroup
		MemberRole string
	}
	var results []result
	err := r.db.Table("run_groups").
		Select("run_groups.*, run_group_members.role as member_role").
		Joins("INNER JOIN run_group_members ON run_group_members.group_id::uuid = run_groups.id").
		Preload("Schedules").
		Where("run_group_members.user_id = ? AND run_group_members.status = ?", userId, "joined").
		Order("run_groups.created_at DESC").
		Find(&results).Error
	if err != nil {
		return nil, nil, err
	}

	groups := make([]entity.RunGroup, len(results))
	roles := make([]string, len(results))
	for i, r := range results {
		groups[i] = r.RunGroup
		roles[i] = r.MemberRole
	}
	return groups, roles, nil
}

func (r *runGroupRepository) GetMemberCount(groupId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.RunGroupMember{}).Where("group_id = ? AND status = ?", groupId, "joined").Count(&count).Error
	return count, err
}
