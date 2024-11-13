package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type RegistrationService interface {
	Registration(registration request.Registration)
}

type RegistrationRepository struct {
	RegistrationRepository repository.UserRepository
	Validate               *validator.Validate
}

func NewRegistrationService(registrationRepo repository.UserRepository, validate *validator.Validate) RegistrationService {
	return &RegistrationRepository{
		RegistrationRepository: registrationRepo,
		Validate:               validate,
	}
}

func (service *RegistrationRepository) Registration(registration request.Registration) {
	err := service.Validate.Struct(registration)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	registrationEntity := entity.Business{
		User: entity.User{
			Email:    registration.Email,
			Password: registration.Password,
			RoleId:   registration.RoleId,
		},
		Name:           registration.Name,
		PhoneNumber:    registration.PhoneNumber,
		OwnerName:      registration.OwnerName,
		Address:        registration.Address,
		BusinessTypeId: registration.BusinessTypeId,
		IsActive:       true,
	}

	service.RegistrationRepository.InsertRegistration((registrationEntity))
}
