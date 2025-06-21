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

	// Email sudah ada?
	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("duplicate email")
	}

	// Buat Business
	business := entity.Business{
		Name:           req.Name,
		OwnerName:      req.OwnerName,
		BusinessTypeId: req.BusinessTypeId,
		Logo:           req.Logo,
		Rating:         req.Rating,
		Image:          req.Image,
		IsActive:       true,
	}
	savedBusiness, err := s.repo.CreateBusiness(business)
	if err != nil {
		return err
	}

	// Buat User
	user := entity.UserBusiness{
		Email:      req.Email,
		Password:   helper.HashAndSalt([]byte(req.Password)),
		RoleId:     req.RoleId,
		BusinessId: savedBusiness.Id,
		IsActive:   true,
		IsVerified: false,
	}
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *registrationService) IsDuplicateEmail(email string) (bool, error) {
	return s.repo.IsEmailExists(email)
}
