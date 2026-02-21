package service

import (
	"math"
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
	repo        repository.RunActivityRepository
	userRepo    repository.UserRepository
	profileRepo repository.RunnerProfileRepository
}

func NewRunActivityService(
	repo repository.RunActivityRepository,
	userRepo repository.UserRepository,
	profileRepo repository.RunnerProfileRepository,
) RunActivityService {
	return &runActivityService{repo: repo, userRepo: userRepo, profileRepo: profileRepo}
}

// Create auto-calculates AvgPace if not provided (pace = duration_min / distance_km).
// After saving, it updates the runner profile AvgPace as a running average.
func (s *runActivityService) Create(userId uuid.UUID, req request.CreateRunActivityRequest) (response.RunActivityDetailResponse, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	// Auto-calculate avg pace if not provided or zero
	avgPace := req.AvgPace
	if avgPace <= 0 && req.Distance > 0 && req.Duration > 0 {
		// pace = minutes per km
		durationMinutes := float64(req.Duration) / 60.0
		avgPace = durationMinutes / req.Distance
		avgPace = math.Round(avgPace*100) / 100
	}

	activity := entity.RunActivity{
		Id:        uuid.New(),
		UserId:    userId,
		Distance:  req.Distance,
		Duration:  req.Duration,
		AvgPace:   avgPace,
		Calories:  req.Calories,
		Source:    req.Source,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&activity); err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	// Update runner profile AvgPace as running average
	s.updateProfileAvgPace(userId, avgPace)

	return s.buildDetailResponse(&activity, user), nil
}

// updateProfileAvgPace recalculates the runner profile average pace
// based on all recorded activities.
func (s *runActivityService) updateProfileAvgPace(userId uuid.UUID, latestPace float64) {
	profile, err := s.profileRepo.FindByUserId(userId)
	if err != nil || profile == nil {
		return
	}

	activities, err := s.repo.FindByUserId(userId)
	if err != nil || len(activities) == 0 {
		return
	}

	var totalPace float64
	var count int
	for _, a := range activities {
		if a.AvgPace > 0 {
			totalPace += a.AvgPace
			count++
		}
	}

	if count > 0 {
		newAvg := math.Round((totalPace/float64(count))*100) / 100
		profile.AvgPace = newAvg
		profile.UpdatedAt = time.Now()
		_ = s.profileRepo.Update(profile)
	}
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

	// Auto-recalculate pace if distance/duration changed but pace not explicitly set
	if req.AvgPace == nil && (req.Distance != nil || req.Duration != nil) {
		if activity.Distance > 0 && activity.Duration > 0 {
			durationMinutes := float64(activity.Duration) / 60.0
			activity.AvgPace = math.Round((durationMinutes/activity.Distance)*100) / 100
		}
	}

	if err := s.repo.Update(activity); err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	// Update profile avg pace
	s.updateProfileAvgPace(activity.UserId, activity.AvgPace)

	user, _ := s.userRepo.FindById(activity.UserId)
	return s.buildDetailResponse(activity, user), nil
}

func (s *runActivityService) FindById(id uuid.UUID) (response.RunActivityDetailResponse, error) {
	activity, err := s.repo.FindById(id)
	if err != nil {
		return response.RunActivityDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(activity.UserId)
	return s.buildDetailResponse(activity, user), nil
}

func (s *runActivityService) FindByUserId(userId uuid.UUID) ([]response.RunActivityDetailResponse, error) {
	activities, err := s.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	user, _ := s.userRepo.FindById(userId)

	var responses []response.RunActivityDetailResponse
	for _, activity := range activities {
		responses = append(responses, s.buildDetailResponse(&activity, user))
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

// -- Response builder --

func (s *runActivityService) buildDetailResponse(activity *entity.RunActivity, user *entity.User) response.RunActivityDetailResponse {
	var userRes *response.UserResponse
	if user != nil {
		userRes = &response.UserResponse{
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
	}
}
