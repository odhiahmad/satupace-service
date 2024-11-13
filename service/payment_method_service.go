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

type PaymentMethodService interface {
	CreatePaymentMethod(paymentMethod request.PaymentMethodCreate)
	UpdatePaymentMethod(paymentMethod request.PaymentMethodUpdate)
	FindById(paymentMethodId int) response.PaymentMethodResponse
	FindAll() []response.PaymentMethodResponse
	Delete(paymentMethodId int)
}

type PaymentMethodRepository struct {
	PaymentMethodRepository repository.PaymentMethodRepository
	Validate                *validator.Validate
}

func NewPaymentMethodService(paymentMethodRepo repository.PaymentMethodRepository, validate *validator.Validate) PaymentMethodService {
	return &PaymentMethodRepository{
		PaymentMethodRepository: paymentMethodRepo,
		Validate:                validate,
	}
}

func (service *PaymentMethodRepository) CreatePaymentMethod(paymentMethod request.PaymentMethodCreate) {
	err := service.Validate.Struct(paymentMethod)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	paymentMethodEntity := entity.PaymentMethod{
		Name: paymentMethod.Name,
	}

	service.PaymentMethodRepository.InsertPaymentMethod((paymentMethodEntity))
}

func (service *PaymentMethodRepository) UpdatePaymentMethod(paymentMethod request.PaymentMethodUpdate) {
	paymentMethodData, err := service.PaymentMethodRepository.FindById(paymentMethod.Id)
	helper.ErrorPanic(err)

	paymentMethodData.Name = paymentMethod.Name

	service.PaymentMethodRepository.UpdatePaymentMethod(paymentMethodData)
}

func (service *PaymentMethodRepository) FindById(paymentMethodId int) response.PaymentMethodResponse {
	paymentMethodData, err := service.PaymentMethodRepository.FindById(paymentMethodId)
	helper.ErrorPanic(err)

	tagResponse := response.PaymentMethodResponse{
		Id:   paymentMethodData.Id,
		Nama: paymentMethodData.Name,
	}
	return tagResponse
}

func (t *PaymentMethodRepository) FindAll() []response.PaymentMethodResponse {
	result := t.PaymentMethodRepository.FindAll()

	var tags []response.PaymentMethodResponse
	for _, value := range result {
		tag := response.PaymentMethodResponse{
			Id:   value.Id,
			Nama: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *PaymentMethodRepository) Delete(paymentMethodId int) {
	t.PaymentMethodRepository.Delete(paymentMethodId)
}
