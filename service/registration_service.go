package service

import (
	"errors"
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
}

type registrationService struct {
	repo       repository.RegistrationRepository
	membership repository.MembershipRepository
	validate   *validator.Validate
}

func NewRegistrationService(repo repository.RegistrationRepository, membership repository.MembershipRepository, validate *validator.Validate) RegistrationService {
	return &registrationService{
		repo:       repo,
		membership: membership,
		validate:   validate,
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

	err = s.repo.CreateMainBranch(&mainBranch)
	if err != nil {
		return response.UserBusinessResponse{}, err
	}

	// 3. Buat User
	user := entity.UserBusiness{
		Email:       req.Email,
		Password:    helper.HashAndSalt([]byte(req.Password)),
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

	// 4. Buat Membership
	startedAt, expiredAt := GetMembershipPeriod(req.Type)
	membership := entity.Membership{
		UserId:    savedUser.Id,
		Type:      req.Type,
		StartDate: startedAt,
		EndDate:   expiredAt,
		IsActive:  true,
	}

	if _, err := s.membership.CreateMembership(membership); err != nil {
		return response.UserBusinessResponse{}, err
	}

	// 5. Kembalikan response user
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
