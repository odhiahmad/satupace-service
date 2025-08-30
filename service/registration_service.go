package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type RegistrationService interface {
	Register(req request.RegistrationRequest) error
	IsDuplicateEmail(email string) (bool, error)
}

type registrationService struct {
	repo        repository.RegistrationRepository
	membership  repository.MembershipRepository
	terminal    repository.TerminalRepository
	validate    *validator.Validate
	redisHelper *helper.RedisHelper
}

func NewRegistrationService(repo repository.RegistrationRepository, membership repository.MembershipRepository, terminal repository.TerminalRepository, validate *validator.Validate, redisHelper *helper.RedisHelper) RegistrationService {
	return &registrationService{
		repo:        repo,
		membership:  membership,
		terminal:    terminal,
		validate:    validate,
		redisHelper: redisHelper,
	}
}

func (s *registrationService) Register(req request.RegistrationRequest) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
		exists, err := s.repo.IsEmailExists(*req.Email)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("email sudah digunakan")
		}
	}

	phoneExists, err := s.repo.IsPhoneNumberExists(req.PhoneNumber)
	if err != nil {
		return err
	}
	if phoneExists {
		return errors.New("nomor telepon sudah digunakan")
	}

	var name = strings.ToLower(req.Name)
	var ownerName = strings.ToLower(req.OwnerName)

	business := entity.Business{
		Name:           name,
		OwnerName:      ownerName,
		BusinessTypeId: &req.BusinessTypeId,
		Image:          req.Image,
		IsActive:       true,
	}

	savedBusiness, err := s.repo.CreateBusiness(business)
	if err != nil {
		return err
	}

	terminalSave := entity.Terminal{
		Id:         uuid.New(),
		BusinessId: savedBusiness.Id,
		Name:       "Terminal 1",
		Location:   "Lokasi Utama",
		IsActive:   helper.BoolPtr(true),
	}

	if _, err := s.terminal.Create(terminalSave); err != nil {
		return err
	}

	hashedPassword := helper.HashAndSalt([]byte(req.Password))
	user := entity.UserBusiness{
		Name:        &ownerName,
		Email:       req.Email,
		Password:    hashedPassword,
		RoleId:      1,
		BusinessId:  savedBusiness.Id,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true,
		IsVerified:  false,
	}

	if _, err := s.repo.CreateUser(user); err != nil {
		return err
	}

	startedAt, expiredAt := GetMembershipPeriod("weekly")
	membership := entity.Membership{
		BusinessId: savedBusiness.Id,
		StartDate:  startedAt,
		EndDate:    expiredAt,
		IsActive:   true,
		Type:       "weekly",
	}

	if _, err := s.membership.CreateMembership(membership); err != nil {
		return err
	}

	if err := s.redisHelper.AllowRequest(req.PhoneNumber, 3, 5*time.Minute); err != nil {
		return err
	}

	otpCode := helper.GenerateOTPCode(6)

	err = s.redisHelper.SaveOTP("otp", req.PhoneNumber, otpCode, 5*time.Minute)
	if err != nil {
		log.Println("Gagal simpan OTP:", err)
		return err
	}

	message := fmt.Sprintf("Kode verifikasi akun kamu adalah: %s", otpCode)
	if err := helper.SendOTPViaWhatsApp(req.PhoneNumber, message); err != nil {
		log.Println("Gagal mengirim OTP WhatsApp:", err)
	}

	return nil
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
