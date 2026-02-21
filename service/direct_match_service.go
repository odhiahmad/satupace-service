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

type DirectMatchService interface {
	// GetCandidates returns compatible runners for the user via MatchingEngine
	GetCandidates(userId uuid.UUID) ([]CandidateResult, error)

	// SendMatchRequest creates a pending match or auto-accepts if reverse match exists
	SendMatchRequest(senderId uuid.UUID, req request.CreateDirectMatchRequest) (response.DirectMatchDetailResponse, error)

	// AcceptMatch accepts a pending match with transaction (update status + create chat room)
	AcceptMatch(matchId uuid.UUID, userId uuid.UUID) (response.DirectMatchDetailResponse, error)

	// RejectMatch rejects a pending match
	RejectMatch(matchId uuid.UUID, userId uuid.UUID) (response.DirectMatchDetailResponse, error)

	FindById(id uuid.UUID) (response.DirectMatchDetailResponse, error)
	FindUserMatches(userId uuid.UUID) ([]response.DirectMatchDetailResponse, error)
	FindMatchesByStatus(userId uuid.UUID, status string) ([]response.DirectMatchDetailResponse, error)
	Delete(id uuid.UUID) error
}

type directMatchService struct {
	repo        repository.DirectMatchRepository
	userRepo    repository.UserRepository
	chatRepo    repository.DirectChatMessageRepository
	profileRepo repository.RunnerProfileRepository
	engine      MatchingEngine
	db          *gorm.DB
}

func NewDirectMatchService(
	repo repository.DirectMatchRepository,
	userRepo repository.UserRepository,
	chatRepo repository.DirectChatMessageRepository,
	profileRepo repository.RunnerProfileRepository,
	engine MatchingEngine,
	db *gorm.DB,
) DirectMatchService {
	return &directMatchService{
		repo:        repo,
		userRepo:    userRepo,
		chatRepo:    chatRepo,
		profileRepo: profileRepo,
		engine:      engine,
		db:          db,
	}
}

// GetCandidates returns compatible runners via the MatchingEngine.
func (s *directMatchService) GetCandidates(userId uuid.UUID) ([]CandidateResult, error) {
	return s.engine.FindDirectCandidates(userId)
}

// SendMatchRequest creates a match request. If the receiver already has a
// pending request TO the sender (reverse match), it auto-accepts both and
// creates a chat room inside a transaction.
func (s *directMatchService) SendMatchRequest(senderId uuid.UUID, req request.CreateDirectMatchRequest) (response.DirectMatchDetailResponse, error) {
	receiverId, err := uuid.Parse(req.User2Id)
	if err != nil {
		return response.DirectMatchDetailResponse{}, errors.New("user_2_id tidak valid")
	}

	if senderId == receiverId {
		return response.DirectMatchDetailResponse{}, errors.New("tidak bisa match dengan diri sendiri")
	}

	// Check both users exist
	sender, err := s.userRepo.FindById(senderId)
	if err != nil {
		return response.DirectMatchDetailResponse{}, errors.New("pengirim tidak ditemukan")
	}
	receiver, err := s.userRepo.FindById(receiverId)
	if err != nil {
		return response.DirectMatchDetailResponse{}, errors.New("penerima tidak ditemukan")
	}

	// Check if match already exists in either direction
	existing, _ := s.repo.FindByUsers(senderId, receiverId)
	if existing != nil {
		return response.DirectMatchDetailResponse{}, errors.New("match sudah ada antara kedua user")
	}

	// Check for reverse match: receiver already sent a pending request to sender
	reverseMatch, _ := s.repo.FindByUsers(receiverId, senderId)
	if reverseMatch != nil && reverseMatch.Status == "pending" {
		// Auto-accept: both like each other! Use transaction
		var result response.DirectMatchDetailResponse
		txErr := s.db.Transaction(func(tx *gorm.DB) error {
			now := time.Now()
			reverseMatch.Status = "accepted"
			reverseMatch.MatchedAt = &now
			if err := tx.Save(reverseMatch).Error; err != nil {
				return err
			}

			// Create initial system chat message
			chatMsg := entity.DirectChatMessage{
				Id:        uuid.New(),
				MatchId:   reverseMatch.Id,
				SenderId:  uuid.Nil, // system message
				Message:   "Match! Kalian saling tertarik. Mulai obrolan sekarang!",
				CreatedAt: now,
			}
			if err := tx.Create(&chatMsg).Error; err != nil {
				return err
			}

			return nil
		})
		if txErr != nil {
			return response.DirectMatchDetailResponse{}, txErr
		}

		result = s.buildMatchResponse(reverseMatch, sender, receiver)
		return result, nil
	}

	// No reverse match - create a new pending match
	match := entity.DirectMatch{
		Id:        uuid.New(),
		User1Id:   senderId,
		User2Id:   receiverId,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&match); err != nil {
		return response.DirectMatchDetailResponse{}, err
	}

	return s.buildMatchResponse(&match, sender, receiver), nil
}

// AcceptMatch accepts a pending match with transaction: update status + create chat room.
func (s *directMatchService) AcceptMatch(matchId uuid.UUID, userId uuid.UUID) (response.DirectMatchDetailResponse, error) {
	match, err := s.repo.FindById(matchId)
	if err != nil {
		return response.DirectMatchDetailResponse{}, errors.New("match tidak ditemukan")
	}

	// Only the receiver (User2) can accept
	if match.User2Id != userId {
		return response.DirectMatchDetailResponse{}, errors.New("hanya penerima yang bisa menerima match")
	}

	if match.Status != "pending" {
		return response.DirectMatchDetailResponse{}, errors.New("match sudah diproses sebelumnya")
	}

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		match.Status = "accepted"
		match.MatchedAt = &now
		if err := tx.Save(match).Error; err != nil {
			return err
		}

		// Create chat room with system message
		chatMsg := entity.DirectChatMessage{
			Id:        uuid.New(),
			MatchId:   match.Id,
			SenderId:  uuid.Nil,
			Message:   "Match diterima! Mulai obrolan sekarang!",
			CreatedAt: now,
		}
		return tx.Create(&chatMsg).Error
	})
	if txErr != nil {
		return response.DirectMatchDetailResponse{}, txErr
	}

	user1, _ := s.userRepo.FindById(match.User1Id)
	user2, _ := s.userRepo.FindById(match.User2Id)
	return s.buildMatchResponse(match, user1, user2), nil
}

// RejectMatch rejects a pending match.
func (s *directMatchService) RejectMatch(matchId uuid.UUID, userId uuid.UUID) (response.DirectMatchDetailResponse, error) {
	match, err := s.repo.FindById(matchId)
	if err != nil {
		return response.DirectMatchDetailResponse{}, errors.New("match tidak ditemukan")
	}

	if match.User2Id != userId {
		return response.DirectMatchDetailResponse{}, errors.New("hanya penerima yang bisa menolak match")
	}

	if match.Status != "pending" {
		return response.DirectMatchDetailResponse{}, errors.New("match sudah diproses sebelumnya")
	}

	match.Status = "rejected"
	if err := s.repo.Update(match); err != nil {
		return response.DirectMatchDetailResponse{}, err
	}

	user1, _ := s.userRepo.FindById(match.User1Id)
	user2, _ := s.userRepo.FindById(match.User2Id)
	return s.buildMatchResponse(match, user1, user2), nil
}

func (s *directMatchService) FindById(id uuid.UUID) (response.DirectMatchDetailResponse, error) {
	match, err := s.repo.FindById(id)
	if err != nil {
		return response.DirectMatchDetailResponse{}, err
	}

	user1, _ := s.userRepo.FindById(match.User1Id)
	user2, _ := s.userRepo.FindById(match.User2Id)
	return s.buildMatchResponse(match, user1, user2), nil
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
		responses = append(responses, s.buildMatchResponse(&match, user1, user2))
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
		responses = append(responses, s.buildMatchResponse(&match, user1, user2))
	}
	return responses, nil
}

func (s *directMatchService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// -- Response builder --

func (s *directMatchService) buildMatchResponse(match *entity.DirectMatch, user1, user2 *entity.User) response.DirectMatchDetailResponse {
	var u1Res, u2Res *response.UserResponse
	if user1 != nil {
		u1Res = &response.UserResponse{
			Id:          user1.Id.String(),
			Name:        user1.Name,
			Email:       user1.Email,
			PhoneNumber: user1.PhoneNumber,
			Gender:      user1.Gender,
			HasProfile:  user1.HasProfile,
			IsVerified:  user1.IsVerified,
			IsActive:    user1.IsActive,
			CreatedAt:   user1.CreatedAt,
			UpdatedAt:   user1.UpdatedAt,
		}
	}
	if user2 != nil {
		u2Res = &response.UserResponse{
			Id:          user2.Id.String(),
			Name:        user2.Name,
			Email:       user2.Email,
			PhoneNumber: user2.PhoneNumber,
			Gender:      user2.Gender,
			HasProfile:  user2.HasProfile,
			IsVerified:  user2.IsVerified,
			IsActive:    user2.IsActive,
			CreatedAt:   user2.CreatedAt,
			UpdatedAt:   user2.UpdatedAt,
		}
	}

	return response.DirectMatchDetailResponse{
		Id:        match.Id.String(),
		User1Id:   match.User1Id.String(),
		User1:     u1Res,
		User2Id:   match.User2Id.String(),
		User2:     u2Res,
		Status:    match.Status,
		CreatedAt: match.CreatedAt,
		MatchedAt: match.MatchedAt,
	}
}
