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

type EmployeeService interface {
	Create(req request.EmployeeRequest) (entity.Employee, error)
	Update(id uuid.UUID, req request.EmployeeUpdateRequest) (entity.Employee, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) response.EmployeeResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.EmployeeResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.EmployeeResponse, string, bool, error)
}

type employeeService struct {
	repo     repository.EmployeeRepository
	validate *validator.Validate
}

func NewEmployeeService(repo repository.EmployeeRepository, validate *validator.Validate) EmployeeService {
	return &employeeService{
		repo:     repo,
		validate: validate,
	}
}

func (s *employeeService) Create(req request.EmployeeRequest) (entity.Employee, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Employee{}, err
	}

	employee := entity.Employee{
		BusinessId:  req.BusinessId,
		Name:        strings.ToLower(req.Name),
		RoleId:      req.RoleId,
		PhoneNumber: req.PhoneNumber,
		Pin:         helper.HashAndSalt([]byte(req.Pin)),
	}

	return s.repo.Create(employee)
}

func (s *employeeService) Update(id uuid.UUID, req request.EmployeeUpdateRequest) (entity.Employee, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Employee{}, err
	}

	employee := entity.Employee{
		Id:          id,
		RoleId:      *req.RoleId,
		Name:        strings.ToLower(*req.Name),
		PhoneNumber: req.PhoneNumber,
		Pin:         helper.HashAndSalt([]byte(*req.Pin)),
	}

	return s.repo.Update(employee)
}

func (s *employeeService) Delete(id uuid.UUID) error {
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

func (s *employeeService) FindById(employeeId uuid.UUID) response.EmployeeResponse {
	employeeData, err := s.repo.FindById(employeeId)
	helper.ErrorPanic(err)

	employeeResponse := mapper.MapEmployee(employeeData)
	return *employeeResponse
}

func (s *employeeService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.EmployeeResponse, int64, error) {
	employees, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.EmployeeResponse
	for _, employee := range employees {
		result = append(result, *mapper.MapEmployee(employee))
	}

	return result, total, nil
}

func (s *employeeService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.EmployeeResponse, string, bool, error) {
	employees, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.EmployeeResponse
	for _, employee := range employees {
		result = append(result, *mapper.MapEmployee(employee))
	}

	return result, nextCursor, hasNext, nil
}
