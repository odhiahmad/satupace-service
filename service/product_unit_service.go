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

type ProductUnitService interface {
	CreateProductUnit(productUnit request.ProductUnitCreate)
	UpdateProductUnit(productUnit request.ProductUnitUpdate)
	FindById(productUnitId int) response.ProductUnitResponse
	FindAll() []response.ProductUnitResponse
	Delete(productUnitId int)
}

type ProductUnitRepository struct {
	ProductUnitRepository repository.ProductUnitRepository
	Validate              *validator.Validate
}

func NewProductUnitService(productUnitRepo repository.ProductUnitRepository, validate *validator.Validate) ProductUnitService {
	return &ProductUnitRepository{
		ProductUnitRepository: productUnitRepo,
		Validate:              validate,
	}
}

func (service *ProductUnitRepository) CreateProductUnit(productUnit request.ProductUnitCreate) {
	err := service.Validate.Struct(productUnit)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	productUnitEntity := entity.ProductUnit{
		Name: productUnit.Name,
	}

	service.ProductUnitRepository.InsertProductUnit((productUnitEntity))
}

func (service *ProductUnitRepository) UpdateProductUnit(productUnit request.ProductUnitUpdate) {
	productUnitData, err := service.ProductUnitRepository.FindById(productUnit.Id)
	helper.ErrorPanic(err)

	productUnitData.Name = productUnit.Name

	service.ProductUnitRepository.UpdateProductUnit(productUnitData)
}

func (service *ProductUnitRepository) FindById(productUnitId int) response.ProductUnitResponse {
	productUnitData, err := service.ProductUnitRepository.FindById(productUnitId)
	helper.ErrorPanic(err)

	tagResponse := response.ProductUnitResponse{
		Id:   productUnitData.Id,
		Nama: productUnitData.Name,
	}
	return tagResponse
}

func (t *ProductUnitRepository) FindAll() []response.ProductUnitResponse {
	result := t.ProductUnitRepository.FindAll()

	var tags []response.ProductUnitResponse
	for _, value := range result {
		tag := response.ProductUnitResponse{
			Id:   value.Id,
			Nama: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *ProductUnitRepository) Delete(productUnitId int) {
	t.ProductUnitRepository.Delete(productUnitId)
}
