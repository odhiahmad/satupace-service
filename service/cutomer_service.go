package service

import (
	"strings"

	"loka-kasir/data/request"
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper"
	"loka-kasir/helper/mapper"
	"loka-kasir/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CustomerService interface {
	Create(req request.CustomerRequest) (response.CustomerResponse, error)
	Update(id uuid.UUID, req request.CustomerRequest) (response.CustomerResponse, error)
	Delete(id uuid.UUID) error
	FindById(roleId uuid.UUID) response.CustomerResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.CustomerResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.CustomerResponse, string, bool, error)
}

type customerService struct {
	repo     repository.CustomerRepository
	validate *validator.Validate
}

func NewCustomerService(repo repository.CustomerRepository, validate *validator.Validate) CustomerService {
	return &customerService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *customerService) Create(req request.CustomerRequest) (response.CustomerResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.CustomerResponse{}, err
	}

	customer := entity.Customer{
		BusinessId: req.BusinessId,
		Name:       strings.ToLower(req.Name),
	}

	createdCustomer, err := s.repo.Create(customer)
	if err != nil {
		return response.CustomerResponse{}, err
	}

	customerResponse := mapper.MapCustomer(&createdCustomer)

	return *customerResponse, nil
}

func (s *customerService) Update(id uuid.UUID, req request.CustomerRequest) (response.CustomerResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.CustomerResponse{}, err
	}

	customer := entity.Customer{
		Id:   id,
		Name: strings.ToLower(req.Name),
	}

	updatedCustomer, err := s.repo.Update(customer)
	if err != nil {
		return response.CustomerResponse{}, err
	}

	customerResponse := mapper.MapCustomer(&updatedCustomer)

	return *customerResponse, nil
}

func (s *customerService) Delete(id uuid.UUID) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	hasRelation, err := s.repo.HasRelation(id)
	if err != nil {
		return err
	}

	var deleteErr error
	if hasRelation {
		deleteErr = s.repo.SoftDelete(id)
	} else {
		deleteErr = s.repo.HardDelete(id)
	}
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (s *customerService) FindById(customerId uuid.UUID) response.CustomerResponse {
	customerData, err := s.repo.FindById(customerId)
	helper.ErrorPanic(err)

	customerResponse := mapper.MapCustomer(&customerData)
	return *customerResponse
}

func (s *customerService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.CustomerResponse, int64, error) {
	customeres, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.CustomerResponse
	for _, customer := range customeres {
		result = append(result, *mapper.MapCustomer(&customer))
	}

	return result, total, nil
}

func (s *customerService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.CustomerResponse, string, bool, error) {
	customers, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.CustomerResponse
	for _, customer := range customers {
		result = append(result, *mapper.MapCustomer(&customer))
	}

	return result, nextCursor, hasNext, nil
}
