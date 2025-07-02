package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type PaymentMethodService interface {
	CreatePaymentMethod(req request.PaymentMethodCreate) (entity.PaymentMethod, error)
	UpdatePaymentMethod(id int, req request.PaymentMethodUpdate) (entity.PaymentMethod, error)
	FindById(id int) (response.PaymentMethodResponse, error)
	FindAll() ([]response.PaymentMethodResponse, error)
	Delete(id int) error
}

type paymentMethodService struct {
	repo     repository.PaymentMethodRepository
	validate *validator.Validate
}

func NewPaymentMethodService(repo repository.PaymentMethodRepository, validate *validator.Validate) PaymentMethodService {
	return &paymentMethodService{
		repo:     repo,
		validate: validate,
	}
}

func (s *paymentMethodService) CreatePaymentMethod(req request.PaymentMethodCreate) (entity.PaymentMethod, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.PaymentMethod{}, err
	}

	method := entity.PaymentMethod{
		Name: req.Name,
	}

	return s.repo.InsertPaymentMethod(method)
}

func (s *paymentMethodService) UpdatePaymentMethod(id int, req request.PaymentMethodUpdate) (entity.PaymentMethod, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.PaymentMethod{}, err
	}

	method, err := s.repo.FindById(id)
	if err != nil {
		return entity.PaymentMethod{}, err
	}

	method.Name = req.Name

	return s.repo.UpdatePaymentMethod(method)
}

func (s *paymentMethodService) FindById(id int) (response.PaymentMethodResponse, error) {
	method, err := s.repo.FindById(id)
	if err != nil {
		return response.PaymentMethodResponse{}, err
	}

	return response.PaymentMethodResponse{
		Id:   method.Id,
		Nama: method.Name,
	}, nil
}

func (s *paymentMethodService) FindAll() ([]response.PaymentMethodResponse, error) {
	methods, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.PaymentMethodResponse
	for _, m := range methods {
		responses = append(responses, response.PaymentMethodResponse{
			Id:   m.Id,
			Nama: m.Name,
		})
	}

	return responses, nil
}

func (s *paymentMethodService) Delete(id int) error {
	return s.repo.Delete(id)
}
