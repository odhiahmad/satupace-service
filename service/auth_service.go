package service

import (
	"fmt"
	"log"

	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
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
		return nil, fmt.Errorf("pengguna tidak ditemukan: %w", err)
	}

	if !comparePassword(user.Password, []byte(password)) {
		return nil, fmt.Errorf("email / nomor HP atau password tidak valid")
	}

	token := service.jwtService.GenerateToken(user.Id)

	role := response.RoleResponse{
		Id:   user.Role.Id,
		Name: user.Role.Name,
	}

	businessType := response.BusinessTypeResponse{
		Id:   user.Business.BusinessType.Id,
		Name: user.Business.BusinessType.Name,
	}

	business := response.BusinessResponse{
		Id:             user.Business.Id,
		Name:           user.Business.Name,
		OwnerName:      user.Business.OwnerName,
		BusinessTypeId: user.Business.BusinessTypeId,
		Image:          user.Business.Image,
		IsActive:       user.Business.IsActive,
		Type:           businessType,
	}

	var branch *response.BusinessBranchResponse
	if user.Branch != nil {
		branch = &response.BusinessBranchResponse{
			Id:          user.Branch.Id,
			PhoneNumber: user.Branch.Phone,
		}
	}

	authRes := response.AuthResponse{
		Id:          user.Id,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Token:       token,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Role:        role,
		Business:    business,
		Branch:      branch,
	}

	return &authRes, nil
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
