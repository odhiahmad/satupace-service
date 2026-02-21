package service

import (
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type SafetyLogService interface {
	Create(userId uuid.UUID, req request.CreateSafetyLogRequest) (response.SafetyLogDetailResponse, error)
	FindById(id uuid.UUID) (response.SafetyLogDetailResponse, error)
	FindByUserId(userId uuid.UUID) ([]response.SafetyLogDetailResponse, error)
	FindByMatchId(matchId uuid.UUID) ([]response.SafetyLogDetailResponse, error)
	FindByStatus(status string) ([]response.SafetyLogDetailResponse, error)
	Delete(id uuid.UUID) error
}

type safetyLogService struct {
	repo     repository.SafetyLogRepository
	userRepo repository.UserRepository
}

func NewSafetyLogService(repo repository.SafetyLogRepository, userRepo repository.UserRepository) SafetyLogService {
	return &safetyLogService{repo: repo, userRepo: userRepo}
}

func (s *safetyLogService) Create(userId uuid.UUID, req request.CreateSafetyLogRequest) (response.SafetyLogDetailResponse, error) {
	matchId, _ := uuid.Parse(req.MatchId)

	log := entity.SafetyLog{
		Id:        uuid.New(),
		UserId:    userId,
		MatchId:   matchId,
		Status:    req.Status,
		Reason:    req.Reason,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&log); err != nil {
		return response.SafetyLogDetailResponse{}, err
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

	return response.SafetyLogDetailResponse{
		Id:        log.Id.String(),
		UserId:    log.UserId.String(),
		User:      userRes,
		MatchId:   log.MatchId.String(),
		Status:    log.Status,
		Reason:    log.Reason,
		CreatedAt: log.CreatedAt,
	}, nil
}

func (s *safetyLogService) FindById(id uuid.UUID) (response.SafetyLogDetailResponse, error) {
	log, err := s.repo.FindById(id)
	if err != nil {
		return response.SafetyLogDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(log.UserId)
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

	return response.SafetyLogDetailResponse{
		Id:        log.Id.String(),
		UserId:    log.UserId.String(),
		User:      userRes,
		MatchId:   log.MatchId.String(),
		Status:    log.Status,
		Reason:    log.Reason,
		CreatedAt: log.CreatedAt,
	}, nil
}

func (s *safetyLogService) FindByUserId(userId uuid.UUID) ([]response.SafetyLogDetailResponse, error) {
	logs, err := s.repo.FindByUserId(userId)
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

	var responses []response.SafetyLogDetailResponse
	for _, log := range logs {
		responses = append(responses, response.SafetyLogDetailResponse{
			Id:        log.Id.String(),
			UserId:    log.UserId.String(),
			User:      userRes,
			MatchId:   log.MatchId.String(),
			Status:    log.Status,
			Reason:    log.Reason,
			CreatedAt: log.CreatedAt,
		})
	}

	return responses, nil
}

func (s *safetyLogService) FindByMatchId(matchId uuid.UUID) ([]response.SafetyLogDetailResponse, error) {
	logs, err := s.repo.FindByMatchId(matchId)
	if err != nil {
		return nil, err
	}

	var responses []response.SafetyLogDetailResponse
	for _, log := range logs {
		user, _ := s.userRepo.FindById(log.UserId)
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

		responses = append(responses, response.SafetyLogDetailResponse{
			Id:        log.Id.String(),
			UserId:    log.UserId.String(),
			User:      userRes,
			MatchId:   log.MatchId.String(),
			Status:    log.Status,
			Reason:    log.Reason,
			CreatedAt: log.CreatedAt,
		})
	}

	return responses, nil
}

func (s *safetyLogService) FindByStatus(status string) ([]response.SafetyLogDetailResponse, error) {
	logs, err := s.repo.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	var responses []response.SafetyLogDetailResponse
	for _, log := range logs {
		user, _ := s.userRepo.FindById(log.UserId)
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

		responses = append(responses, response.SafetyLogDetailResponse{
			Id:        log.Id.String(),
			UserId:    log.UserId.String(),
			User:      userRes,
			MatchId:   log.MatchId.String(),
			Status:    log.Status,
			Reason:    log.Reason,
			CreatedAt: log.CreatedAt,
		})
	}

	return responses, nil
}

func (s *safetyLogService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
