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
	UpdateRole(requesterId uuid.UUID, memberId uuid.UUID, req request.UpdateMemberRoleRequest) (response.RunGroupMemberDetailResponse, error)
	FindById(id uuid.UUID) (response.RunGroupMemberDetailResponse, error)
	FindByGroupId(groupId uuid.UUID) ([]response.RunGroupMemberDetailResponse, error)
	FindByUserId(userId uuid.UUID) ([]response.RunGroupMemberDetailResponse, error)
	Delete(id uuid.UUID) error
	JoinGroup(userId uuid.UUID, groupId uuid.UUID) (response.RunGroupMemberDetailResponse, error)
	LeaveGroup(userId uuid.UUID, groupId uuid.UUID) error
	KickMember(requesterId uuid.UUID, memberId uuid.UUID) error
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
		Role:     "member",
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

// UpdateRole allows owner or admin to change a member's role.
// Only owner can promote/demote to admin. Admins cannot change other admins or owner.
func (s *runGroupMemberService) UpdateRole(requesterId uuid.UUID, memberId uuid.UUID, req request.UpdateMemberRoleRequest) (response.RunGroupMemberDetailResponse, error) {
	targetMember, err := s.repo.FindById(memberId)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("anggota tidak ditemukan")
	}

	requester, err := s.repo.FindByGroupAndUser(targetMember.GroupId, requesterId)
	if err != nil {
		return response.RunGroupMemberDetailResponse{}, errors.New("anda bukan anggota grup ini")
	}

	// Only owner can change roles
	if requester.Role != "owner" {
		return response.RunGroupMemberDetailResponse{}, errors.New("hanya owner yang dapat mengubah role anggota")
	}

	// Cannot change own role
	if targetMember.UserId == requesterId {
		return response.RunGroupMemberDetailResponse{}, errors.New("tidak dapat mengubah role sendiri")
	}

	// Cannot change owner role
	if targetMember.Role == "owner" {
		return response.RunGroupMemberDetailResponse{}, errors.New("tidak dapat mengubah role owner")
	}

	targetMember.Role = req.Role
	if err := s.repo.Update(targetMember); err != nil {
		return response.RunGroupMemberDetailResponse{}, err
	}

	user, _ := s.userRepo.FindById(targetMember.UserId)
	return s.buildResponse(targetMember, user), nil
}

// LeaveGroup allows a member to leave. Owner cannot leave (must delete group instead).
func (s *runGroupMemberService) LeaveGroup(userId uuid.UUID, groupId uuid.UUID) error {
	member, err := s.repo.FindByGroupAndUser(groupId, userId)
	if err != nil {
		return errors.New("anda bukan anggota grup ini")
	}

	if member.Role == "owner" {
		return errors.New("owner tidak dapat meninggalkan grup, hapus grup atau transfer ownership terlebih dahulu")
	}

	if err := s.repo.Delete(member.Id); err != nil {
		return err
	}

	// Re-open group if it was full
	group, _ := s.groupRepo.FindById(groupId)
	if group != nil && group.Status == "full" {
		group.Status = "open"
		s.groupRepo.Update(group)
	}

	return nil
}

// KickMember allows owner or admin to remove a member.
// Admin cannot kick other admins or owner.
func (s *runGroupMemberService) KickMember(requesterId uuid.UUID, memberId uuid.UUID) error {
	targetMember, err := s.repo.FindById(memberId)
	if err != nil {
		return errors.New("anggota tidak ditemukan")
	}

	requester, err := s.repo.FindByGroupAndUser(targetMember.GroupId, requesterId)
	if err != nil {
		return errors.New("anda bukan anggota grup ini")
	}

	// Must be owner or admin
	if requester.Role != "owner" && requester.Role != "admin" {
		return errors.New("hanya owner atau admin yang dapat mengeluarkan anggota")
	}

	// Cannot kick yourself
	if targetMember.UserId == requesterId {
		return errors.New("tidak dapat mengeluarkan diri sendiri, gunakan leave")
	}

	// Cannot kick owner
	if targetMember.Role == "owner" {
		return errors.New("tidak dapat mengeluarkan owner")
	}

	// Admin cannot kick other admins
	if requester.Role == "admin" && targetMember.Role == "admin" {
		return errors.New("admin tidak dapat mengeluarkan admin lain")
	}

	if err := s.repo.Delete(targetMember.Id); err != nil {
		return err
	}

	// Re-open group if it was full
	group, _ := s.groupRepo.FindById(targetMember.GroupId)
	if group != nil && group.Status == "full" {
		group.Status = "open"
		s.groupRepo.Update(group)
	}

	return nil
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
			Role:     "member",
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
			HasProfile:  user.HasProfile,
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
		Role:     member.Role,
		Status:   member.Status,
		JoinedAt: member.JoinedAt,
	}
}
