package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	VerifyCredentialBusiness(identifier string, password string) (*response.AuthResponse, error)
	VerifyOTPToken(req request.VerifyOTPRequest) (*response.AuthResponse, error)
	RetryOTP(req request.RetryOTPRequest) error
	RequestForgotPassword(req request.ForgotPasswordRequest) error
	ResetPassword(req request.ResetPasswordRequest) error
}

type authService struct {
	userRepository         repository.UserRepository
	userBusinessRepository repository.UserBusinessRepository
	jwtService             JWTService
	redisHelper            *helper.RedisHelper
	emailHelper            *helper.EmailHelper
}

func NewAuthService(userRep repository.UserRepository, userBusinessRepository repository.UserBusinessRepository, jwtSvc JWTService, redisHelper *helper.RedisHelper, emailHelper *helper.EmailHelper) AuthService {
	return &authService{
		userRepository:         userRep,
		userBusinessRepository: userBusinessRepository,
		jwtService:             jwtSvc,
		redisHelper:            redisHelper,
		emailHelper:            emailHelper,
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

	if !comparePassword(user.Password, []byte(password)) {
		return nil, helper.ErrInvalidPassword
	}

	if !user.IsVerified {
		otpCode := helper.GenerateOTPCode(6)

		err := service.redisHelper.SaveOTP("otp", user.PhoneNumber, otpCode, 5*time.Minute)
		if err != nil {
			log.Println("Gagal menyimpan OTP di Redis:", err)
		}

		message := fmt.Sprintf("Kode verifikasi akun kamu adalah: %s", otpCode)
		if err := helper.SendOTPViaWhatsApp(user.PhoneNumber, message); err != nil {
			log.Println("Gagal mengirim OTP WhatsApp:", err)
		}

		return nil, helper.ErrEmailNotVerified
	}

	now := time.Now()
	hasActiveMembership := user.Membership.IsActive && user.Membership.EndDate.After(now)

	if !hasActiveMembership {
		return nil, helper.ErrMembershipInactive
	}

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

func (s *authService) VerifyOTPToken(req request.VerifyOTPRequest) (*response.AuthResponse, error) {
	savedOTP, err := s.redisHelper.GetOTP("otp", req.Identifier)
	if err != nil {
		return nil, errors.New("OTP tidak ditemukan atau sudah kedaluwarsa")
	}

	if savedOTP != req.Token {
		return nil, errors.New("OTP tidak valid")
	}

	user, err := s.userBusinessRepository.FindByEmailOrPhone(req.Identifier)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if user.PendingEmail != nil && req.Identifier == *user.PendingEmail {
		user.Email = user.PendingEmail
		user.PendingEmail = nil
		err := s.userBusinessRepository.Update(&user)
		if err != nil {
			return nil, errors.New("gagal memperbarui email")
		}
		_ = s.redisHelper.DeleteOTP("otp", req.Identifier)
	} else {
		if user.IsVerified {
			return nil, errors.New("akun sudah terverifikasi")
		}
		user.IsVerified = true
		err = s.userBusinessRepository.Update(&user)
		if err != nil {
			return nil, errors.New("gagal memverifikasi akun")
		}
		_ = s.redisHelper.DeleteOTP("otp", req.Identifier)
	}

	jwtToken := s.jwtService.GenerateToken(user)
	res := helper.MapAuthResponse(&user, jwtToken)
	return res, nil
}

func (s *authService) RetryOTP(req request.RetryOTPRequest) error {
	if err := s.redisHelper.AllowRequest("retry:"+req.Identifier, 3, 5*time.Minute); err != nil {
		return err
	}

	_ = s.redisHelper.DeleteOTP("otp", req.Identifier)

	newOTP := helper.GenerateOTPCode(6)

	if err := s.redisHelper.SaveOTP("otp", req.Identifier, newOTP, 5*time.Minute); err != nil {
		return errors.New("gagal menyimpan OTP baru")
	}

	if helper.IsEmail(req.Identifier) {
		subject, text, html := helper.BuildVerificationEmail(req.Identifier, newOTP)
		if err := s.emailHelper.Send(req.Identifier, subject, text, html); err != nil {
			return errors.New("gagal mengirim ulang OTP ke email")
		}
	} else {
		message := fmt.Sprintf("Kode verifikasi akun kamu adalah: %s", newOTP)
		if err := helper.SendOTPViaWhatsApp(req.Identifier, message); err != nil {
			return errors.New("gagal mengirim ulang OTP ke WhatsApp")
		}
	}

	return nil
}

func (s *authService) RequestForgotPassword(req request.ForgotPasswordRequest) error {
	user, err := s.userBusinessRepository.FindByEmailOrPhone(req.Identifier)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if err := s.redisHelper.AllowRequest("forgot:"+req.Identifier, 3, 5*time.Minute); err != nil {
		return err
	}

	otpCode := helper.GenerateOTPCode(6)
	err = s.redisHelper.SaveOTP("otp", req.Identifier, otpCode, 5*time.Minute)
	if err != nil {
		return err
	}

	if helper.IsEmail(req.Identifier) {
		subject, text, html := helper.BuildPasswordResetEmail(*user.Email, otpCode)
		if err := s.emailHelper.Send(req.Identifier, subject, text, html); err != nil {
			return err
		}
	} else {
		message := fmt.Sprintf("Kode reset password kamu adalah: %s", otpCode)
		if err := helper.SendOTPViaWhatsApp(req.Identifier, message); err != nil {
			return err
		}
	}

	return nil
}

func (s *authService) ResetPassword(req request.ResetPasswordRequest) error {
	savedOTP, err := s.redisHelper.GetOTP("otp", req.Identifier)
	if err != nil {
		return errors.New("OTP tidak ditemukan atau sudah kedaluwarsa")
	}

	if savedOTP != req.OTP {
		return errors.New("OTP tidak valid")
	}

	user, err := s.userBusinessRepository.FindByEmailOrPhone(req.Identifier)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	hashedPassword := helper.HashAndSalt([]byte(req.NewPassword))
	user.Password = hashedPassword

	if err := s.userBusinessRepository.Update(&user); err != nil {
		return errors.New("gagal mengubah password")
	}

	_ = s.redisHelper.DeleteOTP("otp", req.Identifier)

	return nil
}
