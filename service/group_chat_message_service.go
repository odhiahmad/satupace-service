package service

import (
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type GroupChatMessageService interface {
	Create(userId uuid.UUID, req request.SendGroupChatMessageRequest) (response.GroupChatMessageDetailResponse, error)
	FindByGroupId(groupId uuid.UUID) ([]response.GroupChatMessageDetailResponse, error)
	FindBySenderId(userId uuid.UUID) ([]response.GroupChatMessageDetailResponse, error)
	Delete(id uuid.UUID) error
}

type groupChatMessageService struct {
	repo     repository.GroupChatMessageRepository
	userRepo repository.UserRepository
}

func NewGroupChatMessageService(repo repository.GroupChatMessageRepository, userRepo repository.UserRepository) GroupChatMessageService {
	return &groupChatMessageService{repo: repo, userRepo: userRepo}
}

func (s *groupChatMessageService) Create(userId uuid.UUID, req request.SendGroupChatMessageRequest) (response.GroupChatMessageDetailResponse, error) {
	groupId, _ := uuid.Parse(req.GroupId)

	message := entity.GroupChatMessage{
		Id:        uuid.New(),
		GroupId:   groupId,
		SenderId:  userId,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&message); err != nil {
		return response.GroupChatMessageDetailResponse{}, err
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

	return response.GroupChatMessageDetailResponse{
		Id:        message.Id.String(),
		GroupId:   message.GroupId.String(),
		SenderId:  message.SenderId.String(),
		Sender:    userRes,
		Message:   message.Message,
		CreatedAt: message.CreatedAt,
	}, nil
}

func (s *groupChatMessageService) FindByGroupId(groupId uuid.UUID) ([]response.GroupChatMessageDetailResponse, error) {
	messages, err := s.repo.FindByGroupId(groupId)
	if err != nil {
		return nil, err
	}

	var responses []response.GroupChatMessageDetailResponse
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

		responses = append(responses, response.GroupChatMessageDetailResponse{
			Id:        msg.Id.String(),
			GroupId:   msg.GroupId.String(),
			SenderId:  msg.SenderId.String(),
			Sender:    userRes,
			Message:   msg.Message,
			CreatedAt: msg.CreatedAt,
		})
	}

	return responses, nil
}

func (s *groupChatMessageService) FindBySenderId(userId uuid.UUID) ([]response.GroupChatMessageDetailResponse, error) {
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

	var responses []response.GroupChatMessageDetailResponse
	for _, msg := range messages {
		responses = append(responses, response.GroupChatMessageDetailResponse{
			Id:        msg.Id.String(),
			GroupId:   msg.GroupId.String(),
			SenderId:  msg.SenderId.String(),
			Sender:    userRes,
			Message:   msg.Message,
			CreatedAt: msg.CreatedAt,
		})
	}

	return responses, nil
}

func (s *groupChatMessageService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
