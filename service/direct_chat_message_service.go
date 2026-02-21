package service

import (
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type DirectChatMessageService interface {
	Create(userId uuid.UUID, req request.SendDirectChatMessageRequest) (response.DirectChatMessageDetailResponse, error)
	FindByMatchId(matchId uuid.UUID) ([]response.DirectChatMessageDetailResponse, error)
	FindBySenderId(userId uuid.UUID) ([]response.DirectChatMessageDetailResponse, error)
	Delete(id uuid.UUID) error
}

type directChatMessageService struct {
	repo     repository.DirectChatMessageRepository
	userRepo repository.UserRepository
}

func NewDirectChatMessageService(repo repository.DirectChatMessageRepository, userRepo repository.UserRepository) DirectChatMessageService {
	return &directChatMessageService{repo: repo, userRepo: userRepo}
}

func (s *directChatMessageService) Create(userId uuid.UUID, req request.SendDirectChatMessageRequest) (response.DirectChatMessageDetailResponse, error) {
	matchId, _ := uuid.Parse(req.MatchId)

	message := entity.DirectChatMessage{
		Id:        uuid.New(),
		MatchId:   matchId,
		SenderId:  userId,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&message); err != nil {
		return response.DirectChatMessageDetailResponse{}, err
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

	return response.DirectChatMessageDetailResponse{
		Id:        message.Id.String(),
		MatchId:   message.MatchId.String(),
		SenderId:  message.SenderId.String(),
		Sender:    userRes,
		Message:   message.Message,
		CreatedAt: message.CreatedAt,
	}, nil
}

func (s *directChatMessageService) FindByMatchId(matchId uuid.UUID) ([]response.DirectChatMessageDetailResponse, error) {
	messages, err := s.repo.FindByMatchId(matchId)
	if err != nil {
		return nil, err
	}

	var responses []response.DirectChatMessageDetailResponse
	for _, msg := range messages {
		user, _ := s.userRepo.FindById(msg.SenderId)
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

		responses = append(responses, response.DirectChatMessageDetailResponse{
			Id:        msg.Id.String(),
			MatchId:   msg.MatchId.String(),
			SenderId:  msg.SenderId.String(),
			Sender:    userRes,
			Message:   msg.Message,
			CreatedAt: msg.CreatedAt,
		})
	}

	return responses, nil
}

func (s *directChatMessageService) FindBySenderId(userId uuid.UUID) ([]response.DirectChatMessageDetailResponse, error) {
	messages, err := s.repo.FindBySenderId(userId)
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

	var responses []response.DirectChatMessageDetailResponse
	for _, msg := range messages {
		responses = append(responses, response.DirectChatMessageDetailResponse{
			Id:        msg.Id.String(),
			MatchId:   msg.MatchId.String(),
			SenderId:  msg.SenderId.String(),
			Sender:    userRes,
			Message:   msg.Message,
			CreatedAt: msg.CreatedAt,
		})
	}

	return responses, nil
}

func (s *directChatMessageService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
