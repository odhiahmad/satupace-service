package service

import (
	"errors"
	"time"

	"loka-kasir/data/request"
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper"
	"loka-kasir/helper/mapper"
	"loka-kasir/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService interface {
	Create(req request.EmployeeRequest) error
	Update(id uuid.UUID, req request.EmployeeUpdateRequest) error
	Delete(id uuid.UUID) error
	FindById(roleId uuid.UUID) response.UserBusinessResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.UserBusinessResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.UserBusinessResponse, string, bool, error)
}

type employeeService struct {
	repo     repository.UserBusinessRepository
	validate *validator.Validate
}

func NewEmployeeService(repo repository.UserBusinessRepository, validate *validator.Validate) EmployeeService {
	return &employeeService{repo: repo, validate: validate}
}

func (s *employeeService) Create(req request.EmployeeRequest) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	_, err := s.repo.FindByPhoneAndBusinessId(req.BusinessId, *req.PhoneNumber)
	if err == nil {
		return errors.New("nomor HP sudah terdaftar")
	}

	var hashedPassword string
	if *req.Password != "" {
		passHash, _ := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		hashedPassword = string(passHash)
	}

	pinHash, _ := bcrypt.GenerateFromPassword([]byte(req.PinCode), bcrypt.DefaultCost)

	employee := entity.UserBusiness{
		Id:          uuid.New(),
		RoleId:      req.RoleId,
		BusinessId:  req.BusinessId,
		Name:        &req.Name,
		Email:       req.Email,
		PhoneNumber: *req.PhoneNumber,
		Password:    hashedPassword,
		PinCode:     string(pinHash),
		IsVerified:  true,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateEmployee(&employee); err != nil {
		return err
	}

	return nil
}

func (s *employeeService) Update(id uuid.UUID, req request.EmployeeUpdateRequest) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	employee, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	employee.RoleId = req.RoleId
	employee.Name = &req.Name
	employee.UpdatedAt = time.Now()

	if err := s.repo.Update(&employee); err != nil {
		return err
	}

	return nil
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

func (s *employeeService) FindById(employeeId uuid.UUID) response.UserBusinessResponse {
	employeeData, err := s.repo.FindById(employeeId)
	helper.ErrorPanic(err)

	employeeResponse := mapper.MapUserBusiness(employeeData)
	return *employeeResponse
}

func (s *employeeService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.UserBusinessResponse, int64, error) {
	employeees, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.UserBusinessResponse
	for _, employee := range employeees {
		result = append(result, *mapper.MapUserBusiness(employee))
	}

	return result, total, nil
}

func (s *employeeService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.UserBusinessResponse, string, bool, error) {
	employees, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.UserBusinessResponse
	for _, employee := range employees {
		result = append(result, *mapper.MapUserBusiness(employee))
	}

	return result, nextCursor, hasNext, nil
}
