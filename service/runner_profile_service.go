package service

import (
	"errors"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/helper"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type RunnerProfileService interface {
	// CreateOrUpdate enforces single-profile: creates if none exists, updates if it does
	CreateOrUpdate(userId uuid.UUID, req request.CreateRunnerProfileRequest) (response.RunnerProfileDetailResponse, error)
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

// CreateOrUpdate enforces one profile per user.
// If user already has a profile, it updates it. Otherwise creates a new one.
func (s *runnerProfileService) CreateOrUpdate(userId uuid.UUID, req request.CreateRunnerProfileRequest) (response.RunnerProfileDetailResponse, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, errors.New("user tidak ditemukan")
	}

	// Validate pace range (3.0 - 12.0 min/km)
	if req.AvgPace < 3.0 || req.AvgPace > 12.0 {
		return response.RunnerProfileDetailResponse{}, errors.New("avg_pace harus antara 3.0 - 12.0 min/km")
	}

	// Check if profile already exists for this user
	existing, _ := s.repo.FindByUserId(userId)
	if existing != nil {
		// Update existing profile
		existing.AvgPace = req.AvgPace
		existing.PreferredDistance = req.PreferredDistance
		existing.PreferredTime = req.PreferredTime
		existing.Latitude = req.Latitude
		existing.Longitude = req.Longitude
		existing.WomenOnlyMode = req.WomenOnlyMode
		if req.Image != nil && *req.Image != "" {
			// Upload to Cloudinary if base64 image provided
			imageUrl, err := helper.UploadBase64ToCloudinary(*req.Image, "run-sync/profiles")
			if err != nil {
				return response.RunnerProfileDetailResponse{}, errors.New("gagal upload gambar profil: " + err.Error())
			}
			existing.Image = &imageUrl
		}
		existing.UpdatedAt = time.Now()

		if err := s.repo.Update(existing); err != nil {
			return response.RunnerProfileDetailResponse{}, err
		}

		return s.buildDetailResponse(existing, user), nil
	}

	// Upload image to Cloudinary if provided
	var imagePtr *string
	if req.Image != nil && *req.Image != "" {
		imageUrl, err := helper.UploadBase64ToCloudinary(*req.Image, "run-sync/profiles")
		if err != nil {
			return response.RunnerProfileDetailResponse{}, errors.New("gagal upload gambar profil: " + err.Error())
		}
		imagePtr = &imageUrl
	}

	// Create new profile
	profile := entity.RunnerProfile{
		Id:                uuid.New(),
		UserId:            userId,
		AvgPace:           req.AvgPace,
		PreferredDistance: req.PreferredDistance,
		PreferredTime:     req.PreferredTime,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		WomenOnlyMode:     req.WomenOnlyMode,
		Image:             imagePtr,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.Create(&profile); err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	// Mark user as having a profile
	user.HasProfile = true
	_ = s.userRepo.Update(user)

	return s.buildDetailResponse(&profile, user), nil
}

func (s *runnerProfileService) Update(id uuid.UUID, req request.UpdateRunnerProfileRequest) (response.RunnerProfileDetailResponse, error) {
	profile, err := s.repo.FindById(id)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	if req.AvgPace != nil {
		if *req.AvgPace < 3.0 || *req.AvgPace > 12.0 {
			return response.RunnerProfileDetailResponse{}, errors.New("avg_pace harus antara 3.0 - 12.0 min/km")
		}
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
	if req.Image != nil && *req.Image != "" {
		// Upload to Cloudinary if base64 image provided
		imageUrl, err := helper.UploadBase64ToCloudinary(*req.Image, "run-sync/profiles")
		if err != nil {
			return response.RunnerProfileDetailResponse{}, errors.New("gagal upload gambar profil: " + err.Error())
		}
		profile.Image = &imageUrl
	}
	profile.UpdatedAt = time.Now()

	if err := s.repo.Update(profile); err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(profile.UserId)
	return s.buildDetailResponse(profile, user), nil
}

func (s *runnerProfileService) FindById(id uuid.UUID) (response.RunnerProfileDetailResponse, error) {
	profile, err := s.repo.FindById(id)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(profile.UserId)
	return s.buildDetailResponse(profile, user), nil
}

func (s *runnerProfileService) FindByUserId(userId uuid.UUID) (response.RunnerProfileDetailResponse, error) {
	profile, err := s.repo.FindByUserId(userId)
	if err != nil {
		return response.RunnerProfileDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(profile.UserId)
	return s.buildDetailResponse(profile, user), nil
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

// -- Response builder --

func (s *runnerProfileService) buildDetailResponse(profile *entity.RunnerProfile, user *entity.User) response.RunnerProfileDetailResponse {
	var userRes *response.UserResponse
	if user != nil {
		userRes = &response.UserResponse{
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
	}
}
