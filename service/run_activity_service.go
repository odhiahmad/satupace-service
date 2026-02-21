package service

import (
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type RunActivityService interface {
	Create(userId uuid.UUID, req request.CreateRunActivityRequest) (response.RunActivityDetailResponse, error)
	Update(id uuid.UUID, req request.UpdateRunActivityRequest) (response.RunActivityDetailResponse, error)
	FindById(id uuid.UUID) (response.RunActivityDetailResponse, error)
	FindByUserId(userId uuid.UUID) ([]response.RunActivityDetailResponse, error)
	FindAll() ([]response.RunActivityResponse, error)
	Delete(id uuid.UUID) error
	GetUserStats(userId uuid.UUID) (map[string]interface{}, error)
}

type runActivityService struct {
	repo     repository.RunActivityRepository
	userRepo repository.UserRepository
}

func NewRunActivityService(repo repository.RunActivityRepository, userRepo repository.UserRepository) RunActivityService {
	return &runActivityService{repo: repo, userRepo: userRepo}
}

func (s *runActivityService) Create(userId uuid.UUID, req request.CreateRunActivityRequest) (response.RunActivityDetailResponse, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	activity := entity.RunActivity{
		Id:        uuid.New(),
		UserId:    userId,
		Distance:  req.Distance,
		Duration:  req.Duration,
		AvgPace:   req.AvgPace,
		Calories:  req.Calories,
		Source:    req.Source,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&activity); err != nil {
		return response.RunActivityDetailResponse{}, err
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

	return response.RunActivityDetailResponse{
		Id:        activity.Id.String(),
		UserId:    activity.UserId.String(),
		User:      userRes,
		Distance:  activity.Distance,
		Duration:  activity.Duration,
		AvgPace:   activity.AvgPace,
		Calories:  activity.Calories,
		Source:    activity.Source,
		CreatedAt: activity.CreatedAt,
	}, nil
}

func (s *runActivityService) Update(id uuid.UUID, req request.UpdateRunActivityRequest) (response.RunActivityDetailResponse, error) {
	activity, err := s.repo.FindById(id)
	if err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	if req.Distance != nil {
		activity.Distance = *req.Distance
	}
	if req.Duration != nil {
		activity.Duration = *req.Duration
	}
	if req.AvgPace != nil {
		activity.AvgPace = *req.AvgPace
	}
	if req.Calories != nil {
		activity.Calories = *req.Calories
	}
	if req.Source != nil {
		activity.Source = *req.Source
	}

	if err := s.repo.Update(activity); err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(activity.UserId)
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

	return response.RunActivityDetailResponse{
		Id:        activity.Id.String(),
		UserId:    activity.UserId.String(),
		User:      userRes,
		Distance:  activity.Distance,
		Duration:  activity.Duration,
		AvgPace:   activity.AvgPace,
		Calories:  activity.Calories,
		Source:    activity.Source,
		CreatedAt: activity.CreatedAt,
	}, nil
}

func (s *runActivityService) FindById(id uuid.UUID) (response.RunActivityDetailResponse, error) {
	activity, err := s.repo.FindById(id)
	if err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(activity.UserId)
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

	return response.RunActivityDetailResponse{
		Id:        activity.Id.String(),
		UserId:    activity.UserId.String(),
		User:      userRes,
		Distance:  activity.Distance,
		Duration:  activity.Duration,
		AvgPace:   activity.AvgPace,
		Calories:  activity.Calories,
		Source:    activity.Source,
		CreatedAt: activity.CreatedAt,
	}, nil
}

func (s *runActivityService) FindByUserId(userId uuid.UUID) ([]response.RunActivityDetailResponse, error) {
	activities, err := s.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	user, _ := s.userRepo.FindById(userId)
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

	var responses []response.RunActivityDetailResponse
	for _, activity := range activities {
		responses = append(responses, response.RunActivityDetailResponse{
			Id:        activity.Id.String(),
			UserId:    activity.UserId.String(),
			User:      userRes,
			Distance:  activity.Distance,
			Duration:  activity.Duration,
			AvgPace:   activity.AvgPace,
			Calories:  activity.Calories,
			Source:    activity.Source,
			CreatedAt: activity.CreatedAt,
		})
	}

	return responses, nil
}

func (s *runActivityService) FindAll() ([]response.RunActivityResponse, error) {
	activities, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.RunActivityResponse
	for _, activity := range activities {
		responses = append(responses, response.RunActivityResponse{
			Id:        activity.Id.String(),
			UserId:    activity.UserId.String(),
			Distance:  activity.Distance,
			Duration:  activity.Duration,
			AvgPace:   activity.AvgPace,
			Calories:  activity.Calories,
			Source:    activity.Source,
			CreatedAt: activity.CreatedAt,
		})
	}

	return responses, nil
}

func (s *runActivityService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *runActivityService) GetUserStats(userId uuid.UUID) (map[string]interface{}, error) {
	return s.repo.GetUserStats(userId)
}
