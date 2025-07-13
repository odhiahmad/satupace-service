package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type UserBusinessService interface {
	FindById(id int) (response.UserBusinessResponse, error)
	ChangePassword(req request.ChangePasswordRequest) error
	ChangeEmail(req request.ChangeEmailRequest) error
	ChangePhone(req request.ChangePhoneRequest) error
}

type userBusinessService struct {
	repo        repository.UserBusinessRepository
	redisHelper *helper.RedisHelper
}

func NewUserBusinessService(repo repository.UserBusinessRepository, redisHelper *helper.RedisHelper) UserBusinessService {
	return &userBusinessService{repo: repo, redisHelper: redisHelper}
}

func (s *userBusinessService) FindById(id int) (response.UserBusinessResponse, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return response.UserBusinessResponse{}, err
	}
	return *helper.MapUserBusinessResponse(user), nil
}

func (s *userBusinessService) ChangePassword(req request.ChangePasswordRequest) error {
	user, err := s.repo.FindById(req.Id)
	if err != nil {
		return err
	}

	// Verifikasi password lama
	if !helper.ComparePassword(user.Password, req.OldPassword) {
		return errors.New("password lama salah")
	}

	// Hash password baru
	hashedPassword := helper.HashAndSalt([]byte(req.NewPassword))
	user.Password = hashedPassword

	return s.repo.Update(&user)
}

func (s *userBusinessService) ChangeEmail(req request.ChangeEmailRequest) error {
	user, err := s.repo.FindById(req.Id)
	if err != nil {
		return err
	}

	user.Email = req.Email
	return s.repo.Update(&user)
}

func (s *userBusinessService) ChangePhone(req request.ChangePhoneRequest) error {
	user, err := s.repo.FindById(req.Id)
	if err != nil {
		return err
	}

	user.PhoneNumber = *req.PhoneNumber
	otpCode := helper.GenerateOTPCode(6) // misal "123456"

	err = s.redisHelper.SaveOTP("whatsapp", *req.PhoneNumber, otpCode, 5*time.Minute)
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

func (s *userBusinessService) Delete(userId int) {
	s.Delete(userId)
}
