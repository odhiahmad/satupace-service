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

type DirectMatchService interface {
	Create(user1Id uuid.UUID, req request.CreateDirectMatchRequest) (response.DirectMatchDetailResponse, error)
	Update(id uuid.UUID, req request.UpdateDirectMatchStatusRequest) (response.DirectMatchDetailResponse, error)
	FindById(id uuid.UUID) (response.DirectMatchDetailResponse, error)
	FindUserMatches(userId uuid.UUID) ([]response.DirectMatchDetailResponse, error)
	FindMatchesByStatus(userId uuid.UUID, status string) ([]response.DirectMatchDetailResponse, error)
	Delete(id uuid.UUID) error
}

type directMatchService struct {
	repo     repository.DirectMatchRepository
	userRepo repository.UserRepository
}

func NewDirectMatchService(repo repository.DirectMatchRepository, userRepo repository.UserRepository) DirectMatchService {
	return &directMatchService{repo: repo, userRepo: userRepo}
}

func (s *directMatchService) Create(user1Id uuid.UUID, req request.CreateDirectMatchRequest) (response.DirectMatchDetailResponse, error) {
	user2Id, _ := uuid.Parse(req.User2Id)

	user1, err := s.userRepo.FindById(user1Id)
	if err != nil {
		return response.DirectMatchDetailResponse{}, errors.New("user 1 tidak ditemukan")
	}

	user2, err := s.userRepo.FindById(user2Id)
	if err != nil {
		return response.DirectMatchDetailResponse{}, errors.New("user 2 tidak ditemukan")
	}

	match := entity.DirectMatch{
		Id:        uuid.New(),
		User1Id:   user1Id,
		User2Id:   user2Id,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&match); err != nil {
		return response.DirectMatchDetailResponse{}, err
	}

	user1Res := &response.UserResponse{
		Id:          user1.Id.String(),
		Name:        user1.Name,
		Email:       user1.Email,
		PhoneNumber: user1.PhoneNumber,
		Gender:      user1.Gender,
		IsVerified:  user1.IsVerified,
		IsActive:    user1.IsActive,
		CreatedAt:   user1.CreatedAt,
		UpdatedAt:   user1.UpdatedAt,
	}

	user2Res := &response.UserResponse{
		Id:          user2.Id.String(),
		Name:        user2.Name,
		Email:       user2.Email,
		PhoneNumber: user2.PhoneNumber,
		Gender:      user2.Gender,
		IsVerified:  user2.IsVerified,
		IsActive:    user2.IsActive,
		CreatedAt:   user2.CreatedAt,
		UpdatedAt:   user2.UpdatedAt,
	}

	return response.DirectMatchDetailResponse{
		Id:        match.Id.String(),
		User1Id:   match.User1Id.String(),
		User1:     user1Res,
		User2Id:   match.User2Id.String(),
		User2:     user2Res,
		Status:    match.Status,
		CreatedAt: match.CreatedAt,
		MatchedAt: match.MatchedAt,
	}, nil
}

func (s *directMatchService) Update(id uuid.UUID, req request.UpdateDirectMatchStatusRequest) (response.DirectMatchDetailResponse, error) {
	match, err := s.repo.FindById(id)
	if err != nil {
		return response.DirectMatchDetailResponse{}, err
	}

	match.Status = req.Status
	if req.Status == "accepted" {
		now := time.Now()
		match.MatchedAt = &now
	}

	if err := s.repo.Update(match); err != nil {
		return response.DirectMatchDetailResponse{}, err
	}

	user1, _ := s.userRepo.FindById(match.User1Id)
	user2, _ := s.userRepo.FindById(match.User2Id)

	user1Res := &response.UserResponse{
		Id:          user1.Id.String(),
		Name:        user1.Name,
		Email:       user1.Email,
		PhoneNumber: user1.PhoneNumber,
		Gender:      user1.Gender,
		IsVerified:  user1.IsVerified,
		IsActive:    user1.IsActive,
		CreatedAt:   user1.CreatedAt,
		UpdatedAt:   user1.UpdatedAt,
	}

	user2Res := &response.UserResponse{
		Id:          user2.Id.String(),
		Name:        user2.Name,
		Email:       user2.Email,
		PhoneNumber: user2.PhoneNumber,
		Gender:      user2.Gender,
		IsVerified:  user2.IsVerified,
		IsActive:    user2.IsActive,
		CreatedAt:   user2.CreatedAt,
		UpdatedAt:   user2.UpdatedAt,
	}

	return response.DirectMatchDetailResponse{
		Id:        match.Id.String(),
		User1Id:   match.User1Id.String(),
		User1:     user1Res,
		User2Id:   match.User2Id.String(),
		User2:     user2Res,
		Status:    match.Status,
		CreatedAt: match.CreatedAt,
		MatchedAt: match.MatchedAt,
	}, nil
}

func (s *directMatchService) FindById(id uuid.UUID) (response.DirectMatchDetailResponse, error) {
	match, err := s.repo.FindById(id)
	if err != nil {
		return response.DirectMatchDetailResponse{}, err
	}

	user1, _ := s.userRepo.FindById(match.User1Id)
	user2, _ := s.userRepo.FindById(match.User2Id)

	user1Res := &response.UserResponse{
		Id:          user1.Id.String(),
		Name:        user1.Name,
		Email:       user1.Email,
		PhoneNumber: user1.PhoneNumber,
		Gender:      user1.Gender,
		IsVerified:  user1.IsVerified,
		IsActive:    user1.IsActive,
		CreatedAt:   user1.CreatedAt,
		UpdatedAt:   user1.UpdatedAt,
	}

	user2Res := &response.UserResponse{
		Id:          user2.Id.String(),
		Name:        user2.Name,
		Email:       user2.Email,
		PhoneNumber: user2.PhoneNumber,
		Gender:      user2.Gender,
		IsVerified:  user2.IsVerified,
		IsActive:    user2.IsActive,
		CreatedAt:   user2.CreatedAt,
		UpdatedAt:   user2.UpdatedAt,
	}

	return response.DirectMatchDetailResponse{
		Id:        match.Id.String(),
		User1Id:   match.User1Id.String(),
		User1:     user1Res,
		User2Id:   match.User2Id.String(),
		User2:     user2Res,
		Status:    match.Status,
		CreatedAt: match.CreatedAt,
		MatchedAt: match.MatchedAt,
	}, nil
}

func (s *directMatchService) FindUserMatches(userId uuid.UUID) ([]response.DirectMatchDetailResponse, error) {
	matches, err := s.repo.FindUserMatches(userId)
	if err != nil {
		return nil, err
	}

	var responses []response.DirectMatchDetailResponse
	for _, match := range matches {
		user1, _ := s.userRepo.FindById(match.User1Id)
		user2, _ := s.userRepo.FindById(match.User2Id)

		user1Res := &response.UserResponse{
			Id:          user1.Id.String(),
			Name:        user1.Name,
			Email:       user1.Email,
			PhoneNumber: user1.PhoneNumber,
			Gender:      user1.Gender,
			IsVerified:  user1.IsVerified,
			IsActive:    user1.IsActive,
			CreatedAt:   user1.CreatedAt,
			UpdatedAt:   user1.UpdatedAt,
		}

		user2Res := &response.UserResponse{
			Id:          user2.Id.String(),
			Name:        user2.Name,
			Email:       user2.Email,
			PhoneNumber: user2.PhoneNumber,
			Gender:      user2.Gender,
			IsVerified:  user2.IsVerified,
			IsActive:    user2.IsActive,
			CreatedAt:   user2.CreatedAt,
			UpdatedAt:   user2.UpdatedAt,
		}

		responses = append(responses, response.DirectMatchDetailResponse{
			Id:        match.Id.String(),
			User1Id:   match.User1Id.String(),
			User1:     user1Res,
			User2Id:   match.User2Id.String(),
			User2:     user2Res,
			Status:    match.Status,
			CreatedAt: match.CreatedAt,
			MatchedAt: match.MatchedAt,
		})
	}

	return responses, nil
}

func (s *directMatchService) FindMatchesByStatus(userId uuid.UUID, status string) ([]response.DirectMatchDetailResponse, error) {
	matches, err := s.repo.FindMatchesByStatus(userId, status)
	if err != nil {
		return nil, err
	}

	var responses []response.DirectMatchDetailResponse
	for _, match := range matches {
		user1, _ := s.userRepo.FindById(match.User1Id)
		user2, _ := s.userRepo.FindById(match.User2Id)

		user1Res := &response.UserResponse{
			Id:          user1.Id.String(),
			Name:        user1.Name,
			Email:       user1.Email,
			PhoneNumber: user1.PhoneNumber,
			Gender:      user1.Gender,
			IsVerified:  user1.IsVerified,
			IsActive:    user1.IsActive,
			CreatedAt:   user1.CreatedAt,
			UpdatedAt:   user1.UpdatedAt,
		}

		user2Res := &response.UserResponse{
			Id:          user2.Id.String(),
			Name:        user2.Name,
			Email:       user2.Email,
			PhoneNumber: user2.PhoneNumber,
			Gender:      user2.Gender,
			IsVerified:  user2.IsVerified,
			IsActive:    user2.IsActive,
			CreatedAt:   user2.CreatedAt,
			UpdatedAt:   user2.UpdatedAt,
		}

		responses = append(responses, response.DirectMatchDetailResponse{
			Id:        match.Id.String(),
			User1Id:   match.User1Id.String(),
			User1:     user1Res,
			User2Id:   match.User2Id.String(),
			User2:     user2Res,
			Status:    match.Status,
			CreatedAt: match.CreatedAt,
			MatchedAt: match.MatchedAt,
		})
	}

	return responses, nil
}

func (s *directMatchService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
