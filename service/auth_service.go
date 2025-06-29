package service

import (
	"fmt"
	"log"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	VerifyCredentialBusiness(identifier string, password string) (*entity.UserBusiness, error)
}

type authService struct {
	userRepository         repository.UserRepository
	userBusinessRepository repository.UserBusinessRepository
}

func NewAuthService(userRep repository.UserRepository, userBusinessRepository repository.UserBusinessRepository) AuthService {
	return &authService{
		userRepository:         userRep,
		userBusinessRepository: userBusinessRepository,
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

func (service *authService) VerifyCredentialBusiness(identifier string, password string) (*entity.UserBusiness, error) {
	user, err := service.userBusinessRepository.FindByEmailOrPhone(identifier)
	if err != nil {
		return nil, err
	}

	if !comparePassword(user.Password, []byte(password)) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
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
