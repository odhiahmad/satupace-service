package service

import (
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
}

type authService struct {
	userRepository         repository.UserRepository
	userBusinessRepository repository.UserBusinessRepository
	jwtService             JWTService
}

func NewAuthService(userRep repository.UserRepository, userBusinessRepository repository.UserBusinessRepository, jwtSvc JWTService) AuthService {
	return &authService{
		userRepository:         userRep,
		userBusinessRepository: userBusinessRepository,
		jwtService:             jwtSvc,
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

	// Cek membership aktif
	now := time.Now()
	hasActiveMembership := false
	for _, m := range user.Memberships {
		if m.IsActive && m.EndDate.After(now) {
			hasActiveMembership = true
			break
		}
	}

	if !hasActiveMembership {
		return nil, helper.ErrMembershipInactive
	}

	// Token dan response
	token := service.jwtService.GenerateToken(user.Id)
	res := helper.ToAuthResponse(&user, token)

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
