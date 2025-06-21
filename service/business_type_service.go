package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type BusinessTypeService interface {
	CreateBusinessType(businessType request.BusinessTypeCreate)
	UpdateBusinessType(businessType request.BusinessTypeUpdate)
	FindById(businessTypeId int) response.BusinessTypeResponse
	FindAll() []response.BusinessTypeResponse
	Delete(businessTypeId int)
}

type BusinessTypeRepository struct {
	BusinessTypeRepository repository.BusinessTypeRepository
	Validate               *validator.Validate
}

func NewBusinessTypeService(businessTypeRepo repository.BusinessTypeRepository, validate *validator.Validate) BusinessTypeService {
	return &BusinessTypeRepository{
		BusinessTypeRepository: businessTypeRepo,
		Validate:               validate,
	}
}

func (service *BusinessTypeRepository) CreateBusinessType(businessType request.BusinessTypeCreate) {
	err := service.Validate.Struct(businessType)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	businessTypeEntity := entity.BusinessType{
		Name: businessType.Name,
	}

	service.BusinessTypeRepository.InsertBusinessType((businessTypeEntity))
}

func (service *BusinessTypeRepository) UpdateBusinessType(businessType request.BusinessTypeUpdate) {
	businessTypeData, err := service.BusinessTypeRepository.FindById(businessType.Id)
	helper.ErrorPanic(err)

	businessTypeData.Name = businessType.Name

	service.BusinessTypeRepository.UpdateBusinessType(businessTypeData)
}

func (service *BusinessTypeRepository) FindById(businessTypeId int) response.BusinessTypeResponse {
	businessTypeData, err := service.BusinessTypeRepository.FindById(businessTypeId)
	helper.ErrorPanic(err)

	tagResponse := response.BusinessTypeResponse{
		Id:   businessTypeData.Id,
		Nama: businessTypeData.Name,
	}
	return tagResponse
}

func (t *BusinessTypeRepository) FindAll() []response.BusinessTypeResponse {
	result := t.BusinessTypeRepository.FindAll()

	var tags []response.BusinessTypeResponse
	for _, value := range result {
		tag := response.BusinessTypeResponse{
			Id:   value.Id,
			Nama: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *BusinessTypeRepository) Delete(businessTypeId int) {
	t.BusinessTypeRepository.Delete(businessTypeId)
}
