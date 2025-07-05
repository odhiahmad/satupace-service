package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type RegistrationService interface {
	Register(req request.RegistrationRequest) (response.UserBusinessResponse, error)
	IsDuplicateEmail(email string) (bool, error)
	VerifyOTPToken(phone string, token string) error
	RetryOTP(phone string) error
}

type registrationService struct {
	repo                   repository.RegistrationRepository
	userBusinessRepository repository.UserBusinessRepository
	membership             repository.MembershipRepository
	validate               *validator.Validate
	emailService           EmailService
	redisHelper            *helper.RedisHelper // ← Tambahkan ini
}

func NewRegistrationService(repo repository.RegistrationRepository, membership repository.MembershipRepository, emailService EmailService, validate *validator.Validate, redisHelper *helper.RedisHelper, userBusinessRepository repository.UserBusinessRepository) RegistrationService {
	return &registrationService{
		repo:                   repo,
		membership:             membership,
		validate:               validate,
		emailService:           emailService,
		redisHelper:            redisHelper,
		userBusinessRepository: userBusinessRepository,
	}
}

func (s *registrationService) Register(req request.RegistrationRequest) (response.UserBusinessResponse, error) {
	// Validasi input
	if err := s.validate.Struct(req); err != nil {
		return response.UserBusinessResponse{}, err
	}

	// Cek duplikat email
	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return response.UserBusinessResponse{}, err
	}
	if exists {
		return response.UserBusinessResponse{}, errors.New("email sudah digunakan")
	}

	// 1. Buat Business
	business := entity.Business{
		Name:           req.Name,
		OwnerName:      req.OwnerName,
		BusinessTypeId: req.BusinessTypeId,
		Image:          req.Image,
		IsActive:       true,
	}
	savedBusiness, err := s.repo.CreateBusiness(business)
	if err != nil {
		return response.UserBusinessResponse{}, err
	}

	// 2. Buat Branch Utama
	mainBranch := entity.BusinessBranch{
		BusinessId:  savedBusiness.Id,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Rating:      req.Rating,
		Provinsi:    req.Provinsi,
		Kota:        req.Kota,
		Kecamatan:   req.Kecamatan,
		PostalCode:  req.PostalCode,
		IsMain:      true,
		IsActive:    true,
	}
	if err := s.repo.CreateMainBranch(&mainBranch); err != nil {
		return response.UserBusinessResponse{}, err
	}

	// 3. Hash password
	hashedPassword := helper.HashAndSalt([]byte(req.Password))

	// 4. Generate verification token (bisa dianggap OTP juga)
	otpCode := helper.GenerateOTPCode(6) // misal "123456"

	// Simpan ke Redis selama 5 menit
	err = s.redisHelper.SaveOTP("whatsapp", *req.PhoneNumber, otpCode, 5*time.Minute)
	if err != nil {
		log.Println("Gagal simpan OTP:", err)
		return response.UserBusinessResponse{}, err
	}

	// 5. Buat User
	user := entity.UserBusiness{
		Email:       req.Email,
		Password:    hashedPassword,
		RoleId:      req.RoleId,
		BusinessId:  savedBusiness.Id,
		BranchId:    &mainBranch.Id,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true,
		IsVerified:  false,
	}
	savedUser, err := s.repo.CreateUser(user)
	if err != nil {
		return response.UserBusinessResponse{}, err
	}

	// 6. Buat Membership
	startedAt, expiredAt := GetMembershipPeriod("weekly")
	membership := entity.Membership{
		UserId:    savedUser.Id,
		StartDate: startedAt,
		EndDate:   expiredAt,
		IsActive:  true,
	}
	if _, err := s.membership.CreateMembership(membership); err != nil {
		return response.UserBusinessResponse{}, err
	}

	// 7. Kirim OTP ke WhatsApp user
	message := fmt.Sprintf("Kode verifikasi akun kamu adalah: %s", otpCode)
	if err := helper.SendOTPViaWhatsApp(*req.PhoneNumber, message); err != nil {
		// Logging tapi tetap lanjut
		log.Println("Gagal mengirim OTP WhatsApp:", err)
	}

	// 8. Kembalikan response user
	return *helper.MapUserBusinessResponse(savedUser), nil
}

func (s *registrationService) IsDuplicateEmail(email string) (bool, error) {
	return s.repo.IsEmailExists(email)
}

func GetMembershipPeriod(membershipType string) (time.Time, time.Time) {
	start := time.Now()
	var end time.Time

	switch strings.ToLower(membershipType) {
	case "weekly":
		end = start.AddDate(0, 0, 7)
	case "monthly":
		end = start.AddDate(0, 1, 0)
	case "yearly":
		end = start.AddDate(1, 0, 0)
	default:
		end = start.AddDate(0, 0, 7)
	}

	return start, end
}

func (s *registrationService) VerifyOTPToken(phone string, token string) error {
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

func (s *registrationService) RetryOTP(phone string) error {
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
