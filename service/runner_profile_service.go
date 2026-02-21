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

type RunnerProfileService interface {
	Create(userId uuid.UUID, req request.CreateRunnerProfileRequest) (response.RunnerProfileDetailResponse, error)
	Update(id uuid.UUID, req request.UpdateRunnerProfileRequest) (response.RunnerProfileDetailResponse, error)
	FindById(id uuid.UUID) (response.RunnerProfileDetailResponse, error)
	FindByUserId(userId uuid.UUID) (response.RunnerProfileDetailResponse, error)
	FindAll() ([]response.RunnerProfileResponse, error)
	Delete(id uuid.UUID) error
}

type runnerProfileService struct {
	repo     repository.RunnerProfileRepository
	userRepo repository.UserRepository
}

func NewRunnerProfileService(repo repository.RunnerProfileRepository, userRepo repository.UserRepository) RunnerProfileService {
	return &runnerProfileService{repo: repo, userRepo: userRepo}
}

func (s *runnerProfileService) Create(userId uuid.UUID, req request.CreateRunnerProfileRequest) (response.RunnerProfileDetailResponse, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, errors.New("user tidak ditemukan")
	}

	profile := entity.RunnerProfile{
		Id:                uuid.New(),
		UserId:            userId,
		AvgPace:           req.AvgPace,
		PreferredDistance: req.PreferredDistance,
		PreferredTime:     req.PreferredTime,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		WomenOnlyMode:     req.WomenOnlyMode,
		Image:             req.Image,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.Create(&profile); err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	userRes := &response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return response.RunnerProfileDetailResponse{
		Id:                profile.Id.String(),
		UserId:            profile.UserId.String(),
		User:              userRes,
		AvgPace:           profile.AvgPace,
		PreferredDistance: profile.PreferredDistance,
		PreferredTime:     profile.PreferredTime,
		Latitude:          profile.Latitude,
		Longitude:         profile.Longitude,
		WomenOnlyMode:     profile.WomenOnlyMode,
		Image:             profile.Image,
		IsActive:          profile.IsActive,
		CreatedAt:         profile.CreatedAt,
		UpdatedAt:         profile.UpdatedAt,
	}, nil
}

func (s *runnerProfileService) Update(id uuid.UUID, req request.UpdateRunnerProfileRequest) (response.RunnerProfileDetailResponse, error) {
	profile, err := s.repo.FindById(id)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	if req.AvgPace != nil {
		profile.AvgPace = *req.AvgPace
	}
	if req.PreferredDistance != nil {
		profile.PreferredDistance = *req.PreferredDistance
	}
	if req.PreferredTime != nil {
		profile.PreferredTime = *req.PreferredTime
	}
	if req.Latitude != nil {
		profile.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		profile.Longitude = *req.Longitude
	}
	if req.WomenOnlyMode != nil {
		profile.WomenOnlyMode = *req.WomenOnlyMode
	}
	if req.Image != nil {
		profile.Image = req.Image
	}
	profile.UpdatedAt = time.Now()

	if err := s.repo.Update(profile); err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(profile.UserId)
	userRes := &response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return response.RunnerProfileDetailResponse{
		Id:                profile.Id.String(),
		UserId:            profile.UserId.String(),
		User:              userRes,
		AvgPace:           profile.AvgPace,
		PreferredDistance: profile.PreferredDistance,
		PreferredTime:     profile.PreferredTime,
		Latitude:          profile.Latitude,
		Longitude:         profile.Longitude,
		WomenOnlyMode:     profile.WomenOnlyMode,
		Image:             profile.Image,
		IsActive:          profile.IsActive,
		CreatedAt:         profile.CreatedAt,
		UpdatedAt:         profile.UpdatedAt,
	}, nil
}

func (s *runnerProfileService) FindById(id uuid.UUID) (response.RunnerProfileDetailResponse, error) {
	profile, err := s.repo.FindById(id)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(profile.UserId)
	userRes := &response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return response.RunnerProfileDetailResponse{
		Id:                profile.Id.String(),
		UserId:            profile.UserId.String(),
		User:              userRes,
		AvgPace:           profile.AvgPace,
		PreferredDistance: profile.PreferredDistance,
		PreferredTime:     profile.PreferredTime,
		Latitude:          profile.Latitude,
		Longitude:         profile.Longitude,
		WomenOnlyMode:     profile.WomenOnlyMode,
		Image:             profile.Image,
		IsActive:          profile.IsActive,
		CreatedAt:         profile.CreatedAt,
		UpdatedAt:         profile.UpdatedAt,
	}, nil
}

func (s *runnerProfileService) FindByUserId(userId uuid.UUID) (response.RunnerProfileDetailResponse, error) {
	profile, err := s.repo.FindByUserId(userId)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(profile.UserId)
	userRes := &response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return response.RunnerProfileDetailResponse{
		Id:                profile.Id.String(),
		UserId:            profile.UserId.String(),
		User:              userRes,
		AvgPace:           profile.AvgPace,
		PreferredDistance: profile.PreferredDistance,
		PreferredTime:     profile.PreferredTime,
		Latitude:          profile.Latitude,
		Longitude:         profile.Longitude,
		WomenOnlyMode:     profile.WomenOnlyMode,
		Image:             profile.Image,
		IsActive:          profile.IsActive,
		CreatedAt:         profile.CreatedAt,
		UpdatedAt:         profile.UpdatedAt,
	}, nil
}

func (s *runnerProfileService) FindAll() ([]response.RunnerProfileResponse, error) {
	profiles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.RunnerProfileResponse
	for _, profile := range profiles {
		responses = append(responses, response.RunnerProfileResponse{
			Id:                profile.Id.String(),
			UserId:            profile.UserId.String(),
			AvgPace:           profile.AvgPace,
			PreferredDistance: profile.PreferredDistance,
			PreferredTime:     profile.PreferredTime,
			Latitude:          profile.Latitude,
			Longitude:         profile.Longitude,
			WomenOnlyMode:     profile.WomenOnlyMode,
			Image:             profile.Image,
			IsActive:          profile.IsActive,
			CreatedAt:         profile.CreatedAt,
			UpdatedAt:         profile.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *runnerProfileService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
