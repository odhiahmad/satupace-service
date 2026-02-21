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

type RunGroupMemberService interface {
	Create(req request.CreateRunGroupMemberRequest) (response.RunGroupMemberDetailResponse, error)
	Update(id uuid.UUID, req request.UpdateRunGroupMemberRequest) (response.RunGroupMemberDetailResponse, error)
	FindById(id uuid.UUID) (response.RunGroupMemberDetailResponse, error)
	FindByGroupId(groupId uuid.UUID) ([]response.RunGroupMemberDetailResponse, error)
	FindByUserId(userId uuid.UUID) ([]response.RunGroupMemberDetailResponse, error)
	Delete(id uuid.UUID) error
	JoinGroup(userId uuid.UUID, groupId uuid.UUID) (response.RunGroupMemberDetailResponse, error)
}

type runGroupMemberService struct {
	repo      repository.RunGroupMemberRepository
	userRepo  repository.UserRepository
	groupRepo repository.RunGroupRepository
}

func NewRunGroupMemberService(
	repo repository.RunGroupMemberRepository,
	userRepo repository.UserRepository,
	groupRepo repository.RunGroupRepository,
) RunGroupMemberService {
	return &runGroupMemberService{repo: repo, userRepo: userRepo, groupRepo: groupRepo}
}

func (s *runGroupMemberService) Create(req request.CreateRunGroupMemberRequest) (response.RunGroupMemberDetailResponse, error) {
	groupId, _ := uuid.Parse(req.GroupId)
	userId, _ := uuid.Parse(req.UserId)

	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("user tidak ditemukan")
	}

	if _, err := s.groupRepo.FindById(groupId); err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("grup tidak ditemukan")
	}

	member := entity.RunGroupMember{
		Id:       uuid.New(),
		GroupId:  groupId,
		UserId:   userId,
		Status:   req.Status,
		JoinedAt: time.Now(),
	}

	if err := s.repo.Create(&member); err != nil {
		return response.RunGroupMemberDetailResponse{}, err
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

	return response.RunGroupMemberDetailResponse{
		Id:       member.Id.String(),
		GroupId:  member.GroupId.String(),
		UserId:   member.UserId.String(),
		User:     userRes,
		Status:   member.Status,
		JoinedAt: member.JoinedAt,
	}, nil
}

func (s *runGroupMemberService) Update(id uuid.UUID, req request.UpdateRunGroupMemberRequest) (response.RunGroupMemberDetailResponse, error) {
	member, err := s.repo.FindById(id)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, err
	}

	if req.Status != nil {
		member.Status = *req.Status
	}

	if err := s.repo.Update(member); err != nil {
		return response.RunGroupMemberDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(member.UserId)
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

	return response.RunGroupMemberDetailResponse{
		Id:       member.Id.String(),
		GroupId:  member.GroupId.String(),
		UserId:   member.UserId.String(),
		User:     userRes,
		Status:   member.Status,
		JoinedAt: member.JoinedAt,
	}, nil
}

func (s *runGroupMemberService) FindById(id uuid.UUID) (response.RunGroupMemberDetailResponse, error) {
	member, err := s.repo.FindById(id)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(member.UserId)
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

	return response.RunGroupMemberDetailResponse{
		Id:       member.Id.String(),
		GroupId:  member.GroupId.String(),
		UserId:   member.UserId.String(),
		User:     userRes,
		Status:   member.Status,
		JoinedAt: member.JoinedAt,
	}, nil
}

func (s *runGroupMemberService) FindByGroupId(groupId uuid.UUID) ([]response.RunGroupMemberDetailResponse, error) {
	members, err := s.repo.FindByGroupId(groupId)
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupMemberDetailResponse
	for _, member := range members {
		user, _ := s.userRepo.FindById(member.UserId)
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

		responses = append(responses, response.RunGroupMemberDetailResponse{
			Id:       member.Id.String(),
			GroupId:  member.GroupId.String(),
			UserId:   member.UserId.String(),
			User:     userRes,
			Status:   member.Status,
			JoinedAt: member.JoinedAt,
		})
	}

	return responses, nil
}

func (s *runGroupMemberService) FindByUserId(userId uuid.UUID) ([]response.RunGroupMemberDetailResponse, error) {
	members, err := s.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupMemberDetailResponse
	for _, member := range members {
		user, _ := s.userRepo.FindById(member.UserId)
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

		responses = append(responses, response.RunGroupMemberDetailResponse{
			Id:       member.Id.String(),
			GroupId:  member.GroupId.String(),
			UserId:   member.UserId.String(),
			User:     userRes,
			Status:   member.Status,
			JoinedAt: member.JoinedAt,
		})
	}

	return responses, nil
}

func (s *runGroupMemberService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *runGroupMemberService) JoinGroup(userId uuid.UUID, groupId uuid.UUID) (response.RunGroupMemberDetailResponse, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("user tidak ditemukan")
	}

	if _, err := s.groupRepo.FindById(groupId); err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("grup tidak ditemukan")
	}

	existing, _ := s.repo.FindByGroupAndUser(groupId, userId)
	if existing != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("user sudah bergabung dengan grup ini")
	}

	member := entity.RunGroupMember{
		Id:       uuid.New(),
		GroupId:  groupId,
		UserId:   userId,
		Status:   "joined",
		JoinedAt: time.Now(),
	}

	if err := s.repo.Create(&member); err != nil {
		return response.RunGroupMemberDetailResponse{}, err
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

	return response.RunGroupMemberDetailResponse{
		Id:       member.Id.String(),
		GroupId:  member.GroupId.String(),
		UserId:   member.UserId.String(),
		User:     userRes,
		Status:   member.Status,
		JoinedAt: member.JoinedAt,
	}, nil
}
