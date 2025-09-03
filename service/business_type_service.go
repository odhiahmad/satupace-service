package service

import (
	"log"

	"loka-kasir/data/request"
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper"
	"loka-kasir/repository"

	"github.com/go-playground/validator/v10"
)

type BusinessTypeService interface {
	CreateBusinessType(businessType request.BusinessTypeCreate)
	UpdateBusinessType(businessTypeId int, businessType request.BusinessTypeUpdate)
	FindById(businessTypeId int) response.BusinessTypeResponse
	FindAll() []response.BusinessTypeResponse
	Delete(businessTypeId int)
}

type businessTypeService struct {
	BusinessTypeRepository repository.BusinessTypeRepository
	Validate               *validator.Validate
}

func NewBusinessTypeService(businessTypeRepo repository.BusinessTypeRepository, validate *validator.Validate) BusinessTypeService {
	return &businessTypeService{
		BusinessTypeRepository: businessTypeRepo,
		Validate:               validate,
	}
}

func (service *businessTypeService) CreateBusinessType(businessType request.BusinessTypeCreate) {
	err := service.Validate.Struct(businessType)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	businessTypeEntity := entity.BusinessType{
		Name: businessType.Name,
	}

	service.BusinessTypeRepository.InsertBusinessType((businessTypeEntity))
}

func (service *businessTypeService) UpdateBusinessType(businessTypeId int, businessType request.BusinessTypeUpdate) {
	businessTypeData, err := service.BusinessTypeRepository.FindById(businessTypeId)
	helper.ErrorPanic(err)

	businessTypeData.Name = businessType.Name

	service.BusinessTypeRepository.UpdateBusinessType(businessTypeData)
}

func (service *businessTypeService) FindById(businessTypeId int) response.BusinessTypeResponse {
	businessTypeData, err := service.BusinessTypeRepository.FindById(businessTypeId)
	helper.ErrorPanic(err)

	tagResponse := response.BusinessTypeResponse{
		Id:   businessTypeData.Id,
		Name: businessTypeData.Name,
	}
	return tagResponse
}

func (t *businessTypeService) FindAll() []response.BusinessTypeResponse {
	result := t.BusinessTypeRepository.FindAll()

	var tags []response.BusinessTypeResponse
	for _, value := range result {
		tag := response.BusinessTypeResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *businessTypeService) Delete(businessTypeId int) {
	t.BusinessTypeRepository.Delete(businessTypeId)
}
