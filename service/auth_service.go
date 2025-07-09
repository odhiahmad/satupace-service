package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	VerifyCredentialBusiness(identifier string, password string) (*response.AuthResponse, error)
	VerifyOTPToken(phone string, token string) error
	RetryOTP(phone string) error
}

type authService struct {
	userRepository         repository.UserRepository
	userBusinessRepository repository.UserBusinessRepository
	jwtService             JWTService
	redisHelper            *helper.RedisHelper
}

func NewAuthService(userRep repository.UserRepository, userBusinessRepository repository.UserBusinessRepository, jwtSvc JWTService, redisHelper *helper.RedisHelper) AuthService {
	return &authService{
		userRepository:         userRep,
		userBusinessRepository: userBusinessRepository,
		jwtService:             jwtSvc,
		redisHelper:            redisHelper,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (service *authService) VerifyCredentialBusiness(identifier string, password string) (*response.AuthResponse, error) {
	user, err := service.userBusinessRepository.FindByEmailOrPhone(identifier)
	if err != nil {
		return nil, helper.ErrUserNotFound
	}

	if !user.IsVerified {
		return nil, helper.ErrEmailNotVerified
	}

	if !comparePassword(user.Password, []byte(password)) {
		return nil, helper.ErrInvalidPassword
	}

	// Cek membership aktif
	now := time.Now()
	hasActiveMembership := false

	if user.Membership.IsActive && user.Membership.EndDate.After(now) {
		hasActiveMembership = true
	}

	if !hasActiveMembership {
		return nil, helper.ErrMembershipInactive
	}

	// Token dan response
	token := service.jwtService.GenerateToken(user)
	res := helper.MapAuthResponse(&user, token)

	return res, nil
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (s *authService) VerifyOTPToken(phone string, token string) error {
	// Ambil OTP dari Redis
	savedOTP, err := s.redisHelper.GetOTP("whatsapp", phone)
	if err != nil {
		return errors.New("OTP tidak ditemukan atau sudah kedaluwarsa")
	}

	// Cocokkan token
	if savedOTP != token {
		return errors.New("OTP tidak valid")
	}

	// Update is_verified = true berdasarkan phone number
	user, err := s.userBusinessRepository.FindByEmailOrPhone(phone)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("akun sudah terverifikasi")
	}

	user.IsVerified = true

	err = s.userBusinessRepository.Update(&user)
	if err != nil {
		return errors.New("gagal memverifikasi akun")
	}

	// (Opsional) Hapus OTP dari Redis setelah verifikasi berhasil
	_ = s.redisHelper.DeleteOTP("whatsapp", phone)

	return nil
}

func (s *authService) RetryOTP(phone string) error {
	// Ambil OTP dari Redis
	otpCode, err := s.redisHelper.GetOTP("whatsapp", phone)
	if err != nil {
		return errors.New("OTP tidak ditemukan atau sudah kedaluwarsa")
	}

	message := fmt.Sprintf("Kode verifikasi akun kamu adalah: %s", otpCode)

	// Kirim ulang OTP dengan retry selama TTL Redis masih aktif
	err = s.redisHelper.RetryUntilRedisKeyExpired(
		"whatsapp",
		phone,
		2*time.Second,
		func() error {
			return helper.SendOTPViaWhatsApp(phone, message)
		},
	)
	if err != nil {
		log.Println("Gagal mengirim ulang OTP:", err)
		return errors.New("gagal mengirim ulang OTP")
	}

	return nil
}
