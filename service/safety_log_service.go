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

const (
	// AutoSuspendThreshold is the number of reports before auto-suspend
	AutoSuspendThreshold = 5
)

type SafetyLogService interface {
	// ReportUser creates a safety report, increments target report count,
	// and auto-suspends the target if threshold is reached
	ReportUser(reporterId uuid.UUID, req request.CreateSafetyLogRequest) (response.SafetyLogDetailResponse, error)
	FindById(id uuid.UUID) (response.SafetyLogDetailResponse, error)
	FindByUserId(userId uuid.UUID) ([]response.SafetyLogDetailResponse, error)
	FindByMatchId(matchId uuid.UUID) ([]response.SafetyLogDetailResponse, error)
	FindByStatus(status string) ([]response.SafetyLogDetailResponse, error)
	Delete(id uuid.UUID) error
}

type safetyLogService struct {
	repo     repository.SafetyLogRepository
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewSafetyLogService(repo repository.SafetyLogRepository, userRepo repository.UserRepository, db *gorm.DB) SafetyLogService {
	return &safetyLogService{repo: repo, userRepo: userRepo, db: db}
}

// ReportUser creates a safety log and handles report counting + auto-suspend.
// req.MatchId is used as the reported user's ID.
func (s *safetyLogService) ReportUser(reporterId uuid.UUID, req request.CreateSafetyLogRequest) (response.SafetyLogDetailResponse, error) {
	reportedUserId, err := uuid.Parse(req.MatchId)
	if err != nil {
		return response.SafetyLogDetailResponse{}, errors.New("match_id (reported user) tidak valid")
	}

	if reporterId == reportedUserId {
		return response.SafetyLogDetailResponse{}, errors.New("tidak bisa melaporkan diri sendiri")
	}

	// Verify reporter exists
	reporter, err := s.userRepo.FindById(reporterId)
	if err != nil {
		return response.SafetyLogDetailResponse{}, errors.New("pelapor tidak ditemukan")
	}

	// Verify reported user exists
	reportedUser, err := s.userRepo.FindById(reportedUserId)
	if err != nil {
		return response.SafetyLogDetailResponse{}, errors.New("user yang dilaporkan tidak ditemukan")
	}

	log := entity.SafetyLog{
		Id:        uuid.New(),
		UserId:    reporterId,
		MatchId:   reportedUserId, // MatchId stores the reported user ID
		Status:    req.Status,
		Reason:    req.Reason,
		CreatedAt: time.Now(),
	}

	// Transaction: create log + increment report count + auto-suspend
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		// Increment report count on reported user
		reportedUser.ReportCount++
		if err := tx.Model(&entity.User{}).Where("id = ?", reportedUserId).
			Update("report_count", reportedUser.ReportCount).Error; err != nil {
			return err
		}

		// Auto-suspend if threshold reached
		if reportedUser.ReportCount >= AutoSuspendThreshold {
			if err := tx.Model(&entity.User{}).Where("id = ?", reportedUserId).
				Updates(map[string]interface{}{
					"is_suspended": true,
					"is_active":    false,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return response.SafetyLogDetailResponse{}, txErr
	}

	reporterRes := s.buildUserResponse(reporter)

	return response.SafetyLogDetailResponse{
		Id:        log.Id.String(),
		UserId:    log.UserId.String(),
		User:      reporterRes,
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
	userRes := s.buildUserResponse(user)

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
	userRes := s.buildUserResponse(user)

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
		userRes := s.buildUserResponse(user)

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
		userRes := s.buildUserResponse(user)

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

// -- Helper --

func (s *safetyLogService) buildUserResponse(user *entity.User) *response.UserResponse {
	if user == nil {
		return nil
	}
	return &response.UserResponse{
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
