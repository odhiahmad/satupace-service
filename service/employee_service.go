package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/helper/mapper"
	"github.com/odhiahmad/kasirku-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService interface {
	Create(req request.EmployeeRequest) (*entity.UserBusiness, error)
	Update(id uuid.UUID, req request.EmployeeRequest) (*entity.UserBusiness, error)
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

func (s *employeeService) Create(req request.EmployeeRequest) (*entity.UserBusiness, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	_, err := s.repo.FindByPhoneAndBusinessId(req.BusinessId, req.PhoneNumber)
	if err == nil {
		return nil, errors.New("nomor HP sudah terdaftar")
	}

	var hashedPassword string
	if req.Password != "" {
		passHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		hashedPassword = string(passHash)
	}

	pinHash, _ := bcrypt.GenerateFromPassword([]byte(req.PinCode), bcrypt.DefaultCost)

	employee := entity.UserBusiness{
		Id:          uuid.New(),
		RoleId:      req.RoleId,
		BusinessId:  req.BusinessId,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		PinCode:     string(pinHash),
		IsVerified:  true,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateEmployee(&employee); err != nil {
		return nil, err
	}

	return &employee, nil
}

func (s *employeeService) Update(id uuid.UUID, req request.EmployeeRequest) (*entity.UserBusiness, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	employee, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	employee.RoleId = req.RoleId
	employee.BusinessId = req.BusinessId
	employee.Email = req.Email
	employee.PhoneNumber = req.PhoneNumber
	employee.UpdatedAt = time.Now()

	if req.Password != "" {
		passHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		employee.Password = string(passHash)
	}

	if req.PinCode != "" {
		pinHash, _ := bcrypt.GenerateFromPassword([]byte(req.PinCode), bcrypt.DefaultCost)
		employee.PinCode = string(pinHash)
	}

	if err := s.repo.Update(&employee); err != nil {
		return nil, err
	}

	return &employee, nil
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
