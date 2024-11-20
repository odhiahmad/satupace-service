package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type RegistrationService interface {
	Registration(registration request.Registration)
	IsDuplicateEmail(user string) bool
}

type registrationService struct {
	businessRepository       repository.BusinessRepository
	businessBranchRepository repository.BusinessBranchRepository
	userBusinessRepository   repository.UserBusinessRepository
	Validate                 *validator.Validate
}

func NewRegistrationService(businessRepository repository.BusinessRepository, businessBranchRepository repository.BusinessBranchRepository, userBusinessRepository repository.UserBusinessRepository, validate *validator.Validate) RegistrationService {
	return &registrationService{
		businessRepository:       businessRepository,
		businessBranchRepository: businessBranchRepository,
		userBusinessRepository:   userBusinessRepository,
		Validate:                 validate,
	}
}

func (service *registrationService) Registration(registration request.Registration) {
	err := service.Validate.Struct(registration)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}

	business := entity.Business{
		Name:           registration.Name,
		OwnerName:      registration.OwnerName,
		BusinessTypeID: registration.BusinessTypeId,
		IsActive:       true,
	}

	businessData := service.businessRepository.InsertBusiness(business)

	log.Println("business", businessData)

	for _, branch := range registration.Branch {
		branches := entity.BusinessBranch{
			BusinessID:  businessData.ID,
			Pic:         branch.Pic,
			PhoneNumber: branch.PhoneNumber,
			Address:     branch.Address,
		}
		service.businessBranchRepository.InsertBusinessBranch(branches)
	}

	userBusiness := entity.UserBusiness{
		Email:      registration.Email,
		Password:   helper.HashAndSalt([]byte(registration.Password)),
		RoleID:     registration.RoleId,
		BusinessID: businessData.ID,
	}

	service.userBusinessRepository.InsertUserBusiness(userBusiness)
}

func (service *registrationService) IsDuplicateEmail(email string) bool {
	res := service.userBusinessRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}
