package service

import (
	"errors"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	db        *gorm.DB
}

func NewRunGroupMemberService(
	repo repository.RunGroupMemberRepository,
	userRepo repository.UserRepository,
	groupRepo repository.RunGroupRepository,
	db *gorm.DB,
) RunGroupMemberService {
	return &runGroupMemberService{repo: repo, userRepo: userRepo, groupRepo: groupRepo, db: db}
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

	return s.buildResponse(&member, user), nil
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
	return s.buildResponse(member, user), nil
}

func (s *runGroupMemberService) FindById(id uuid.UUID) (response.RunGroupMemberDetailResponse, error) {
	member, err := s.repo.FindById(id)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(member.UserId)
	return s.buildResponse(member, user), nil
}

func (s *runGroupMemberService) FindByGroupId(groupId uuid.UUID) ([]response.RunGroupMemberDetailResponse, error) {
	members, err := s.repo.FindByGroupId(groupId)
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupMemberDetailResponse
	for _, member := range members {
		user, _ := s.userRepo.FindById(member.UserId)
		responses = append(responses, s.buildResponse(&member, user))
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
		responses = append(responses, s.buildResponse(&member, user))
	}
	return responses, nil
}

func (s *runGroupMemberService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// JoinGroup uses a GORM transaction to:
// 1. Validate user + group exist
// 2. Check womenOnly gender restriction
// 3. Check if already a member
// 4. Check if group is full
// 5. Insert member
// 6. Auto-set group status to "full" if needed
func (s *runGroupMemberService) JoinGroup(userId uuid.UUID, groupId uuid.UUID) (response.RunGroupMemberDetailResponse, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("user tidak ditemukan")
	}

	group, err := s.groupRepo.FindById(groupId)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("grup tidak ditemukan")
	}

	if group.Status != "open" {
		return response.RunGroupMemberDetailResponse{}, errors.New("grup sudah penuh atau sudah selesai")
	}

	// Women-only gender check
	if group.IsWomenOnly {
		if user.Gender == nil || *user.Gender != "female" {
			return response.RunGroupMemberDetailResponse{}, errors.New("grup ini hanya untuk perempuan")
		}
	}

	existing, _ := s.repo.FindByGroupAndUser(groupId, userId)
	if existing != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("user sudah bergabung dengan grup ini")
	}

	memberCount, _ := s.groupRepo.GetMemberCount(groupId)
	if group.MaxMember > 0 && int(memberCount) >= group.MaxMember {
		return response.RunGroupMemberDetailResponse{}, errors.New("grup sudah penuh")
	}

	var member entity.RunGroupMember

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		member = entity.RunGroupMember{
			Id:       uuid.New(),
			GroupId:  groupId,
			UserId:   userId,
			Status:   "joined",
			JoinedAt: time.Now(),
		}

		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		newCount := memberCount + 1
		if group.MaxMember > 0 && int(newCount) >= group.MaxMember {
			if err := tx.Model(&entity.RunGroup{}).Where("id = ?", groupId).Update("status", "full").Error; err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return response.RunGroupMemberDetailResponse{}, txErr
	}

	return s.buildResponse(&member, user), nil
}

func (s *runGroupMemberService) buildResponse(
	member *entity.RunGroupMember,
	user *entity.User,
) response.RunGroupMemberDetailResponse {
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

	return response.RunGroupMemberDetailResponse{
		Id:       member.Id.String(),
		GroupId:  member.GroupId.String(),
		UserId:   member.UserId.String(),
		User:     userRes,
		Status:   member.Status,
		JoinedAt: member.JoinedAt,
	}
}
