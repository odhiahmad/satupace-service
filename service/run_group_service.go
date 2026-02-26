package service

import (
	"errors"
	"math"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type RunGroupService interface {
	Create(createdBy uuid.UUID, req request.CreateRunGroupRequest) (response.RunGroupDetailResponse, error)
	Update(id uuid.UUID, req request.UpdateRunGroupRequest) (response.RunGroupDetailResponse, error)
	FindById(id uuid.UUID) (response.RunGroupDetailResponse, error)
	FindAll(filter request.RunGroupFilterRequest) ([]response.RunGroupResponse, error)
	FindByStatus(status string) ([]response.RunGroupResponse, error)
	Delete(id uuid.UUID) error
	FindByCreatedBy(userId uuid.UUID) ([]response.RunGroupResponse, error)
	FindMyGroups(userId uuid.UUID) ([]response.RunGroupResponse, error)
}

type runGroupService struct {
	repo       repository.RunGroupRepository
	userRepo   repository.UserRepository
	memberRepo repository.RunGroupMemberRepository
}

func NewRunGroupService(repo repository.RunGroupRepository, userRepo repository.UserRepository, memberRepo repository.RunGroupMemberRepository) RunGroupService {
	return &runGroupService{repo: repo, userRepo: userRepo, memberRepo: memberRepo}
}

func (s *runGroupService) Create(createdBy uuid.UUID, req request.CreateRunGroupRequest) (response.RunGroupDetailResponse, error) {
	user, err := s.userRepo.FindById(createdBy)
	if err != nil {
		return response.RunGroupDetailResponse{}, errors.New("user tidak ditemukan")
	}

	scheduledAt, _ := time.Parse(time.RFC3339, req.ScheduledAt)

	group := entity.RunGroup{
		Id:                uuid.New(),
		Name:              req.Name,
		MinPace:           req.MinPace,
		MaxPace:           req.MaxPace,
		PreferredDistance: req.PreferredDistance,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		MeetingPoint:      req.MeetingPoint,
		ScheduledAt:       scheduledAt,
		MaxMember:         req.MaxMember,
		IsWomenOnly:       req.IsWomenOnly,
		Status:            "open",
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
	}

	if err := s.repo.Create(&group); err != nil {
		return response.RunGroupDetailResponse{}, err
	}

	// Auto-add creator as owner member
	ownerMember := entity.RunGroupMember{
		Id:       uuid.New(),
		GroupId:  group.Id,
		UserId:   createdBy,
		Role:     "owner",
		Status:   "joined",
		JoinedAt: time.Now(),
	}
	s.memberRepo.Create(&ownerMember)

	creatorRes := &response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return response.RunGroupDetailResponse{
		Id:                group.Id.String(),
		Name:              group.Name,
		MinPace:           group.MinPace,
		MaxPace:           group.MaxPace,
		PreferredDistance: group.PreferredDistance,
		Latitude:          group.Latitude,
		Longitude:         group.Longitude,
		MeetingPoint:      group.MeetingPoint,
		ScheduledAt:       group.ScheduledAt,
		MaxMember:         group.MaxMember,
		IsWomenOnly:       group.IsWomenOnly,
		Status:            group.Status,
		CreatedBy:         group.CreatedBy.String(),
		Creator:           creatorRes,
		MemberCount:       1, // owner just joined
		Schedules:         []response.RunGroupScheduleResponse{},
		CreatedAt:         group.CreatedAt,
	}, nil
}

func (s *runGroupService) Update(id uuid.UUID, req request.UpdateRunGroupRequest) (response.RunGroupDetailResponse, error) {
	group, err := s.repo.FindById(id)
	if err != nil {
		return response.RunGroupDetailResponse{}, err
	}

	if req.Name != nil {
		group.Name = req.Name
	}
	if req.MinPace != nil {
		group.MinPace = *req.MinPace
	}
	if req.MaxPace != nil {
		group.MaxPace = *req.MaxPace
	}
	if req.PreferredDistance != nil {
		group.PreferredDistance = *req.PreferredDistance
	}
	if req.Latitude != nil {
		group.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		group.Longitude = *req.Longitude
	}
	if req.MeetingPoint != nil {
		group.MeetingPoint = *req.MeetingPoint
	}
	if req.ScheduledAt != nil {
		scheduledAt, _ := time.Parse(time.RFC3339, *req.ScheduledAt)
		group.ScheduledAt = scheduledAt
	}
	if req.MaxMember != nil {
		group.MaxMember = *req.MaxMember
	}
	if req.IsWomenOnly != nil {
		group.IsWomenOnly = *req.IsWomenOnly
	}
	if req.Status != nil {
		group.Status = *req.Status
	}

	if err := s.repo.Update(group); err != nil {
		return response.RunGroupDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(group.CreatedBy)
	var creatorRes *response.UserResponse
	if user != nil {
		creatorRes = &response.UserResponse{
			Id:          user.Id.String(),
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Gender:      user.Gender,
			HasProfile:  user.HasProfile,
			IsVerified:  user.IsVerified,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}
	}

	memberCount, _ := s.repo.GetMemberCount(group.Id)

	return response.RunGroupDetailResponse{
		Id:                group.Id.String(),
		Name:              group.Name,
		MinPace:           group.MinPace,
		MaxPace:           group.MaxPace,
		PreferredDistance: group.PreferredDistance,
		Latitude:          group.Latitude,
		Longitude:         group.Longitude,
		MeetingPoint:      group.MeetingPoint,
		ScheduledAt:       group.ScheduledAt,
		MaxMember:         group.MaxMember,
		IsWomenOnly:       group.IsWomenOnly,
		Status:            group.Status,
		CreatedBy:         group.CreatedBy.String(),
		Creator:           creatorRes,
		MemberCount:       int(memberCount),
		Schedules:         mapGroupSchedules(group.Schedules),
		CreatedAt:         group.CreatedAt,
	}, nil
}

func (s *runGroupService) FindById(id uuid.UUID) (response.RunGroupDetailResponse, error) {
	group, err := s.repo.FindById(id)
	if err != nil {
		return response.RunGroupDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(group.CreatedBy)
	var creatorRes *response.UserResponse
	if user != nil {
		creatorRes = &response.UserResponse{
			Id:          user.Id.String(),
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Gender:      user.Gender,
			HasProfile:  user.HasProfile,
			IsVerified:  user.IsVerified,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}
	}

	memberCount, _ := s.repo.GetMemberCount(group.Id)

	return response.RunGroupDetailResponse{
		Id:                group.Id.String(),
		Name:              group.Name,
		MinPace:           group.MinPace,
		MaxPace:           group.MaxPace,
		PreferredDistance: group.PreferredDistance,
		Latitude:          group.Latitude,
		Longitude:         group.Longitude,
		MeetingPoint:      group.MeetingPoint,
		ScheduledAt:       group.ScheduledAt,
		MaxMember:         group.MaxMember,
		IsWomenOnly:       group.IsWomenOnly,
		Status:            group.Status,
		CreatedBy:         group.CreatedBy.String(),
		Creator:           creatorRes,
		MemberCount:       int(memberCount),
		Schedules:         mapGroupSchedules(group.Schedules),
		CreatedAt:         group.CreatedAt,
	}, nil
}

func (s *runGroupService) FindAll(filter request.RunGroupFilterRequest) ([]response.RunGroupResponse, error) {
	groups, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	// Parse user location for distance calculation
	var userLat, userLng float64
	hasLocation := false
	if filter.Latitude != "" && filter.Longitude != "" {
		lat, errLat := strconv.ParseFloat(filter.Latitude, 64)
		lng, errLng := strconv.ParseFloat(filter.Longitude, 64)
		if errLat == nil && errLng == nil {
			userLat, userLng = lat, lng
			hasLocation = true
		}
	}

	var responses []response.RunGroupResponse
	for _, group := range groups {
		memberCount, _ := s.repo.GetMemberCount(group.Id)
		res := response.RunGroupResponse{
			Id:                group.Id.String(),
			Name:              group.Name,
			MinPace:           group.MinPace,
			MaxPace:           group.MaxPace,
			PreferredDistance: group.PreferredDistance,
			Latitude:          group.Latitude,
			Longitude:         group.Longitude,
			MeetingPoint:      group.MeetingPoint,
			ScheduledAt:       group.ScheduledAt,
			MaxMember:         group.MaxMember,
			MemberCount:       int(memberCount),
			IsWomenOnly:       group.IsWomenOnly,
			Status:            group.Status,
			CreatedBy:         group.CreatedBy.String(),
			CreatedAt:         group.CreatedAt,
		}

		if hasLocation {
			dist := haversineDistance(userLat, userLng, group.Latitude, group.Longitude)
			dist = math.Round(dist*100) / 100
			res.DistanceKm = &dist
		}

		responses = append(responses, res)
	}

	return responses, nil
}

func (s *runGroupService) FindByStatus(status string) ([]response.RunGroupResponse, error) {
	groups, err := s.repo.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupResponse
	for _, group := range groups {
		memberCount, _ := s.repo.GetMemberCount(group.Id)
		responses = append(responses, response.RunGroupResponse{
			Id:                group.Id.String(),
			Name:              group.Name,
			MinPace:           group.MinPace,
			MaxPace:           group.MaxPace,
			PreferredDistance: group.PreferredDistance,
			Latitude:          group.Latitude,
			Longitude:         group.Longitude,
			MeetingPoint:      group.MeetingPoint,
			ScheduledAt:       group.ScheduledAt,
			MaxMember:         group.MaxMember,
			MemberCount:       int(memberCount),
			IsWomenOnly:       group.IsWomenOnly,
			Status:            group.Status,
			CreatedBy:         group.CreatedBy.String(),
			CreatedAt:         group.CreatedAt,
		})
	}

	return responses, nil
}

func (s *runGroupService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *runGroupService) FindByCreatedBy(userId uuid.UUID) ([]response.RunGroupResponse, error) {
	groups, err := s.repo.FindByCreatedBy(userId)
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupResponse
	for _, group := range groups {
		memberCount, _ := s.repo.GetMemberCount(group.Id)
		responses = append(responses, response.RunGroupResponse{
			Id:                group.Id.String(),
			Name:              group.Name,
			MinPace:           group.MinPace,
			MaxPace:           group.MaxPace,
			PreferredDistance: group.PreferredDistance,
			Latitude:          group.Latitude,
			Longitude:         group.Longitude,
			MeetingPoint:      group.MeetingPoint,
			ScheduledAt:       group.ScheduledAt,
			MaxMember:         group.MaxMember,
			MemberCount:       int(memberCount),
			IsWomenOnly:       group.IsWomenOnly,
			Status:            group.Status,
			CreatedBy:         group.CreatedBy.String(),
			CreatedAt:         group.CreatedAt,
		})
	}

	return responses, nil
}

func (s *runGroupService) FindMyGroups(userId uuid.UUID) ([]response.RunGroupResponse, error) {
	groups, roles, err := s.repo.FindByMembership(userId)
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupResponse
	for i, group := range groups {
		memberCount, _ := s.repo.GetMemberCount(group.Id)
		responses = append(responses, response.RunGroupResponse{
			Id:                group.Id.String(),
			Name:              group.Name,
			MinPace:           group.MinPace,
			MaxPace:           group.MaxPace,
			PreferredDistance: group.PreferredDistance,
			Latitude:          group.Latitude,
			Longitude:         group.Longitude,
			MeetingPoint:      group.MeetingPoint,
			ScheduledAt:       group.ScheduledAt,
			MaxMember:         group.MaxMember,
			MemberCount:       int(memberCount),
			IsWomenOnly:       group.IsWomenOnly,
			Status:            group.Status,
			CreatedBy:         group.CreatedBy.String(),
			MyRole:            roles[i],
			CreatedAt:         group.CreatedAt,
		})
	}

	return responses, nil
}

// haversineDistance calculates the distance (km) between two lat/lng points.
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371.0 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

var groupDayNames = [7]string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

// mapGroupSchedules converts a slice of RunGroupSchedule entities to response DTOs.
func mapGroupSchedules(schedules []*entity.RunGroupSchedule) []response.RunGroupScheduleResponse {
	result := make([]response.RunGroupScheduleResponse, 0, len(schedules))
	for _, s := range schedules {
		dayName := ""
		if s.DayOfWeek >= 0 && s.DayOfWeek <= 6 {
			dayName = groupDayNames[s.DayOfWeek]
		}
		result = append(result, response.RunGroupScheduleResponse{
			Id:        s.Id.String(),
			GroupId:   s.GroupId.String(),
			DayOfWeek: s.DayOfWeek,
			DayName:   dayName,
			StartTime: s.StartTime,
			IsActive:  s.IsActive,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		})
	}
	return result
}
