package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/helper"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type BiometricService interface {
	// Registration flow
	RegisterStart(userId uuid.UUID, req request.BiometricRegisterStartRequest) (response.BiometricChallengeResponse, error)
	RegisterFinish(userId uuid.UUID, req request.BiometricRegisterFinishRequest) (response.BiometricCredentialResponse, error)

	// Authentication flow
	LoginStart(req request.BiometricLoginStartRequest) (response.BiometricChallengeResponse, error)
	LoginFinish(req request.BiometricLoginFinishRequest) (response.UserResponse, string, string, error) // returns user + access token + refresh token

	// Credential management
	GetCredentials(userId uuid.UUID) ([]response.BiometricCredentialResponse, error)
	DeleteCredential(userId uuid.UUID, credentialId uuid.UUID) error
}

type biometricService struct {
	repo       repository.BiometricRepository
	userRepo   repository.UserRepository
	jwtService JWTService
	redis      *helper.RedisHelper
}

func NewBiometricService(
	repo repository.BiometricRepository,
	userRepo repository.UserRepository,
	jwtService JWTService,
	redis *helper.RedisHelper,
) BiometricService {
	return &biometricService{
		repo:       repo,
		userRepo:   userRepo,
		jwtService: jwtService,
		redis:      redis,
	}
}

// generateChallenge creates a cryptographically secure random challenge
func generateChallenge() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("gagal membuat challenge: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// verifySignature verifies the HMAC-SHA256 signature of the challenge using the public key
func verifySignature(publicKeyHex, challenge, signatureHex string) bool {
	pubKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, pubKeyBytes)
	mac.Write([]byte(challenge))
	expectedSig := mac.Sum(nil)

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedSig, sigBytes)
}

// RegisterStart generates a challenge for biometric registration
func (s *biometricService) RegisterStart(userId uuid.UUID, req request.BiometricRegisterStartRequest) (response.BiometricChallengeResponse, error) {
	// Verify user exists
	_, err := s.userRepo.FindById(userId)
	if err != nil {
		return response.BiometricChallengeResponse{}, errors.New("user tidak ditemukan")
	}

	// Generate challenge
	challenge, err := generateChallenge()
	if err != nil {
		return response.BiometricChallengeResponse{}, err
	}

	// Store challenge in Redis (5 minutes expiry)
	if err := s.redis.SaveOTP("biometric_register", userId.String(), challenge, 5*time.Minute); err != nil {
		return response.BiometricChallengeResponse{}, errors.New("gagal menyimpan challenge")
	}

	return response.BiometricChallengeResponse{
		Challenge: challenge,
	}, nil
}

// RegisterFinish verifies the signature and stores the biometric credential
func (s *biometricService) RegisterFinish(userId uuid.UUID, req request.BiometricRegisterFinishRequest) (response.BiometricCredentialResponse, error) {
	// Get stored challenge from Redis
	storedChallenge, err := s.redis.GetOTP("biometric_register", userId.String())
	if err != nil {
		return response.BiometricCredentialResponse{}, errors.New("challenge tidak valid atau sudah kadaluarsa")
	}

	// Verify challenge matches
	if storedChallenge != req.Challenge {
		return response.BiometricCredentialResponse{}, errors.New("challenge tidak cocok")
	}

	// Verify signature using the provided public key
	if !verifySignature(req.PublicKey, req.Challenge, req.Signature) {
		return response.BiometricCredentialResponse{}, errors.New("signature tidak valid")
	}

	// Check if credential already exists
	existing, _ := s.repo.FindByCredentialId(req.CredentialId)
	if existing != nil {
		return response.BiometricCredentialResponse{}, errors.New("credential sudah terdaftar")
	}

	// Save biometric credential
	now := time.Now()
	biometric := entity.UserBiometric{
		Id:           uuid.New(),
		UserId:       userId,
		CredentialId: req.CredentialId,
		PublicKey:    req.PublicKey,
		DeviceName:   req.DeviceName,
		IsActive:     true,
		LastUsedAt:   &now,
		CreatedAt:    now,
	}

	if err := s.repo.Create(&biometric); err != nil {
		return response.BiometricCredentialResponse{}, errors.New("gagal menyimpan credential biometrik")
	}

	// Delete challenge from Redis
	s.redis.DeleteOTP("biometric_register", userId.String())

	return response.BiometricCredentialResponse{
		Id:           biometric.Id.String(),
		UserId:       biometric.UserId.String(),
		CredentialId: biometric.CredentialId,
		DeviceName:   biometric.DeviceName,
		IsActive:     biometric.IsActive,
		LastUsedAt:   biometric.LastUsedAt,
		CreatedAt:    biometric.CreatedAt,
	}, nil
}

// LoginStart generates a challenge for biometric login
func (s *biometricService) LoginStart(req request.BiometricLoginStartRequest) (response.BiometricChallengeResponse, error) {
	// Find user by email or phone
	user, err := s.userRepo.FindByEmailOrPhone(req.Identifier)
	if err != nil {
		return response.BiometricChallengeResponse{}, errors.New("user tidak ditemukan")
	}

	// Check user has biometric credentials
	credentials, err := s.repo.FindByUserId(user.Id)
	if err != nil || len(credentials) == 0 {
		return response.BiometricChallengeResponse{}, errors.New("biometrik belum didaftarkan untuk akun ini")
	}

	// Generate challenge
	challenge, err := generateChallenge()
	if err != nil {
		return response.BiometricChallengeResponse{}, err
	}

	// Store challenge in Redis (5 minutes expiry)
	if err := s.redis.SaveOTP("biometric_login", user.Id.String(), challenge, 5*time.Minute); err != nil {
		return response.BiometricChallengeResponse{}, errors.New("gagal menyimpan challenge")
	}

	return response.BiometricChallengeResponse{
		Challenge: challenge,
	}, nil
}

// LoginFinish verifies the biometric signature and returns a JWT token
func (s *biometricService) LoginFinish(req request.BiometricLoginFinishRequest) (response.UserResponse, string, string, error) {
	// Find biometric credential
	biometric, err := s.repo.FindByCredentialId(req.CredentialId)
	if err != nil {
		return response.UserResponse{}, "", "", errors.New("credential biometrik tidak ditemukan")
	}

	// Get stored challenge from Redis
	storedChallenge, err := s.redis.GetOTP("biometric_login", biometric.UserId.String())
	if err != nil {
		return response.UserResponse{}, "", "", errors.New("challenge tidak valid atau sudah kadaluarsa")
	}

	// Verify challenge matches
	if storedChallenge != req.Challenge {
		return response.UserResponse{}, "", "", errors.New("challenge tidak cocok")
	}

	// Verify signature
	if !verifySignature(biometric.PublicKey, req.Challenge, req.Signature) {
		return response.UserResponse{}, "", "", errors.New("signature biometrik tidak valid")
	}

	// Find user
	user, err := s.userRepo.FindById(biometric.UserId)
	if err != nil {
		return response.UserResponse{}, "", "", errors.New("user tidak ditemukan")
	}

	// Check user status
	if !user.IsVerified {
		return response.UserResponse{}, "", "", errors.New("akun belum diverifikasi")
	}
	if !user.IsActive {
		return response.UserResponse{}, "", "", errors.New("akun tidak aktif")
	}
	if user.IsSuspended {
		return response.UserResponse{}, "", "", errors.New("akun Anda telah disuspend")
	}

	// Update last used
	now := time.Now()
	biometric.LastUsedAt = &now
	_ = s.repo.Update(biometric)

	// Delete challenge from Redis
	s.redis.DeleteOTP("biometric_login", biometric.UserId.String())

	// Generate JWT tokens
	expiryTime := time.Now().Add(1 * time.Hour)
	accessToken := s.jwtService.GenerateToken(user.Id.String(), user.PhoneNumber, user.Email, expiryTime)
	refreshToken := s.jwtService.GenerateRefreshToken(user.Id.String())

	userRes := response.UserResponse{
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

	return userRes, accessToken, refreshToken, nil
}

// GetCredentials returns all biometric credentials for a user
func (s *biometricService) GetCredentials(userId uuid.UUID) ([]response.BiometricCredentialResponse, error) {
	credentials, err := s.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	var res []response.BiometricCredentialResponse
	for _, c := range credentials {
		res = append(res, response.BiometricCredentialResponse{
			Id:           c.Id.String(),
			UserId:       c.UserId.String(),
			CredentialId: c.CredentialId,
			DeviceName:   c.DeviceName,
			IsActive:     c.IsActive,
			LastUsedAt:   c.LastUsedAt,
			CreatedAt:    c.CreatedAt,
		})
	}

	return res, nil
}

// DeleteCredential removes a biometric credential
func (s *biometricService) DeleteCredential(userId uuid.UUID, credentialId uuid.UUID) error {
	// Verify credential belongs to user
	credentials, err := s.repo.FindByUserId(userId)
	if err != nil {
		return errors.New("credential tidak ditemukan")
	}

	for _, c := range credentials {
		if c.Id == credentialId {
			return s.repo.Delete(credentialId)
		}
	}

	return errors.New("credential tidak dimiliki oleh user ini")
}
