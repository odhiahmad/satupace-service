package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/helper/mapper"
	"github.com/odhiahmad/kasirku-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserBusinessService interface {
	FindById(id uuid.UUID) (response.UserBusinessResponse, error)
	ChangePassword(req request.ChangePasswordRequest) error
	ChangeEmail(req request.ChangeEmailRequest) error
	ChangePhone(req request.ChangePhoneRequest) error
	CreateEmployee(req request.CreateEmployeeRequest) (*entity.UserBusiness, error)
}

type userBusinessService struct {
	repo        repository.UserBusinessRepository
	redisHelper *helper.RedisHelper
	emailHelper *helper.EmailHelper
}

func NewUserBusinessService(repo repository.UserBusinessRepository, redisHelper *helper.RedisHelper, emailHelper *helper.EmailHelper) UserBusinessService {
	return &userBusinessService{repo: repo, redisHelper: redisHelper, emailHelper: emailHelper}
}

func (s *userBusinessService) CreateEmployee(req request.CreateEmployeeRequest) (*entity.UserBusiness, error) {
	existing, _ := s.repo.FindByPhoneAndBusinessId(req.BusinessId, req.PhoneNumber)
	if existing != nil {
		return nil, errors.New("nomor HP sudah terdaftar")
	}

	var hashedPassword string
	if req.Password != "" {
		passHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		hashedPassword = string(passHash)
	}

	pinHash, _ := bcrypt.GenerateFromPassword([]byte(req.PinCode), bcrypt.DefaultCost)

	user := entity.UserBusiness{
		Id:          uuid.New(),
		RoleId:      req.RoleId,
		BusinessId:  req.BusinessId,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		PinCode:     string(pinHash),
		IsVerified:  true,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateEmployee(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userBusinessService) FindById(id uuid.UUID) (response.UserBusinessResponse, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return response.UserBusinessResponse{}, err
	}
	return *mapper.MapUserBusiness(user), nil
}

func (s *userBusinessService) ChangePassword(req request.ChangePasswordRequest) error {
	user, err := s.repo.FindById(req.Id)
	if err != nil {
		return err
	}

	if !helper.ComparePassword(user.Password, req.OldPassword) {
		return errors.New("password lama salah")
	}

	hashedPassword := helper.HashAndSalt([]byte(req.NewPassword))
	user.Password = hashedPassword

	return s.repo.Update(&user)
}

func (s *userBusinessService) ChangeEmail(req request.ChangeEmailRequest) error {
	user, err := s.repo.FindById(req.Id)
	if err != nil {
		return err
	}

	otpCode := helper.GenerateOTPCode(6)

	err = s.redisHelper.SaveOTP("otp", *req.Email, otpCode, 5*time.Minute)
	if err != nil {
		log.Println("gagal menyimpan otp verifikasi email baru di redis:", err)
		return err
	}

	subject, text, html := helper.BuildLinkEmailVerification(*req.Email, otpCode)
	if err := s.emailHelper.Send(*req.Email, subject, text, html); err != nil {
		log.Println("gagal mengirim otp verifikasi email baru:", err)
		return err
	}

	user.PendingEmail = req.Email
	return s.repo.Update(&user)
}

func (s *userBusinessService) ChangePhone(req request.ChangePhoneRequest) error {
	user, err := s.repo.FindById(req.Id)
	if err != nil {
		return err
	}

	if err := s.redisHelper.AllowRequest("retry:"+*req.PhoneNumber, 3, 5*time.Minute); err != nil {
		return err
	}

	user.PhoneNumber = *req.PhoneNumber
	otpCode := helper.GenerateOTPCode(6)

	err = s.redisHelper.SaveOTP("otp", *req.PhoneNumber, otpCode, 5*time.Minute)
	if err != nil {
		log.Println("Gagal simpan OTP:", err)
		return err
	}

	message := fmt.Sprintf("Kode verifikasi akun kamu adalah: %s", otpCode)
	if err := helper.SendOTPViaWhatsApp(*req.PhoneNumber, message); err != nil {
		log.Println("Gagal mengirim OTP WhatsApp:", err)
	}

	return s.repo.Update(&user)
}

func (s *userBusinessService) Delete(userId uuid.UUID) {
	s.repo.Delete(userId)
}
