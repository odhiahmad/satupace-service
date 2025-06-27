package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
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
	repo     repository.RegistrationRepository
	validate *validator.Validate
}

func NewRegistrationService(repo repository.RegistrationRepository, validate *validator.Validate) RegistrationService {
	return &registrationService{
		repo:     repo,
		validate: validate,
	}
}

func (s *registrationService) Register(req request.RegistrationRequest) error {
	// Validasi input
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	// Cek duplikat email
	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email sudah digunakan")
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
		return err
	}

	// 2. Buat Branch Utama (default)
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

	err = s.repo.CreateMainBranch(&mainBranch) // butuh method tambahan di repository
	if err != nil {
		return err
	}

	// 3. Buat User (dengan cabang utama)
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

	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *registrationService) IsDuplicateEmail(email string) (bool, error) {
	return s.repo.IsEmailExists(email)
}
