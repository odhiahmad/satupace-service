package service

import (
	"errors"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type RunGroupService interface {
	Create(createdBy uuid.UUID, req request.CreateRunGroupRequest) (response.RunGroupDetailResponse, error)
	Update(id uuid.UUID, req request.UpdateRunGroupRequest) (response.RunGroupDetailResponse, error)
	FindById(id uuid.UUID) (response.RunGroupDetailResponse, error)
	FindAll() ([]response.RunGroupResponse, error)
	FindByStatus(status string) ([]response.RunGroupResponse, error)
	Delete(id uuid.UUID) error
	FindByCreatedBy(userId uuid.UUID) ([]response.RunGroupResponse, error)
}

type runGroupService struct {
	repo     repository.RunGroupRepository
	userRepo repository.UserRepository
}

func NewRunGroupService(repo repository.RunGroupRepository, userRepo repository.UserRepository) RunGroupService {
	return &runGroupService{repo: repo, userRepo: userRepo}
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
		AvgPace:           req.AvgPace,
		PreferredDistance: req.PreferredDistance,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
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
		AvgPace:           group.AvgPace,
		PreferredDistance: group.PreferredDistance,
		Latitude:          group.Latitude,
		Longitude:         group.Longitude,
		ScheduledAt:       group.ScheduledAt,
		MaxMember:         group.MaxMember,
		IsWomenOnly:       group.IsWomenOnly,
		Status:            group.Status,
		CreatedBy:         group.CreatedBy.String(),
		Creator:           creatorRes,
		MemberCount:       0,
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
	if req.AvgPace != nil {
		group.AvgPace = *req.AvgPace
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
		AvgPace:           group.AvgPace,
		PreferredDistance: group.PreferredDistance,
		Latitude:          group.Latitude,
		Longitude:         group.Longitude,
		ScheduledAt:       group.ScheduledAt,
		MaxMember:         group.MaxMember,
		IsWomenOnly:       group.IsWomenOnly,
		Status:            group.Status,
		CreatedBy:         group.CreatedBy.String(),
		Creator:           creatorRes,
		MemberCount:       int(memberCount),
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
		AvgPace:           group.AvgPace,
		PreferredDistance: group.PreferredDistance,
		Latitude:          group.Latitude,
		Longitude:         group.Longitude,
		ScheduledAt:       group.ScheduledAt,
		MaxMember:         group.MaxMember,
		IsWomenOnly:       group.IsWomenOnly,
		Status:            group.Status,
		CreatedBy:         group.CreatedBy.String(),
		Creator:           creatorRes,
		MemberCount:       int(memberCount),
		CreatedAt:         group.CreatedAt,
	}, nil
}

func (s *runGroupService) FindAll() ([]response.RunGroupResponse, error) {
	groups, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupResponse
	for _, group := range groups {
		responses = append(responses, response.RunGroupResponse{
			Id:                group.Id.String(),
			Name:              group.Name,
			AvgPace:           group.AvgPace,
			PreferredDistance: group.PreferredDistance,
			Latitude:          group.Latitude,
			Longitude:         group.Longitude,
			ScheduledAt:       group.ScheduledAt,
			MaxMember:         group.MaxMember,
			IsWomenOnly:       group.IsWomenOnly,
			Status:            group.Status,
			CreatedBy:         group.CreatedBy.String(),
			CreatedAt:         group.CreatedAt,
		})
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
		responses = append(responses, response.RunGroupResponse{
			Id:                group.Id.String(),
			Name:              group.Name,
			AvgPace:           group.AvgPace,
			PreferredDistance: group.PreferredDistance,
			Latitude:          group.Latitude,
			Longitude:         group.Longitude,
			ScheduledAt:       group.ScheduledAt,
			MaxMember:         group.MaxMember,
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
		responses = append(responses, response.RunGroupResponse{
			Id:                group.Id.String(),
			Name:              group.Name,
			AvgPace:           group.AvgPace,
			PreferredDistance: group.PreferredDistance,
			Latitude:          group.Latitude,
			Longitude:         group.Longitude,
			ScheduledAt:       group.ScheduledAt,
			MaxMember:         group.MaxMember,
			IsWomenOnly:       group.IsWomenOnly,
			Status:            group.Status,
			CreatedBy:         group.CreatedBy.String(),
			CreatedAt:         group.CreatedAt,
		})
	}

	return responses, nil
}
