package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type OrderTypeService interface {
	Create(req request.OrderTypeRequest) (response.OrderTypeResponse, error)
	Update(id uuid.UUID, req request.OrderTypeRequest) (response.OrderTypeResponse, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) response.OrderTypeResponse
	FindWithPagination(pagination request.Pagination) ([]response.OrderTypeResponse, int64, error)
	FindWithPaginationCursor(pagination request.Pagination) ([]response.OrderTypeResponse, string, bool, error)
}

type orderTypeService struct {
	repo     repository.OrderTypeRepository
	validate *validator.Validate
}

func NewOrderTypeService(repo repository.OrderTypeRepository, validate *validator.Validate) OrderTypeService {
	return &orderTypeService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *orderTypeService) Create(req request.OrderTypeRequest) (response.OrderTypeResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.OrderTypeResponse{}, err
	}

	orderType := entity.OrderType{
		Name: strings.ToLower(req.Name),
	}

	createdOrderType, err := s.repo.Create(orderType)
	if err != nil {
		return response.OrderTypeResponse{}, err
	}

	orderTypeResponse := helper.MapOrderType(&createdOrderType)

	return *orderTypeResponse, nil
}

func (s *orderTypeService) Update(id uuid.UUID, req request.OrderTypeRequest) (response.OrderTypeResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.OrderTypeResponse{}, err
	}

	orderType := entity.OrderType{
		Id:   id,
		Name: strings.ToLower(req.Name),
	}

	updatedOrderType, err := s.repo.Update(orderType)
	if err != nil {
		return response.OrderTypeResponse{}, err
	}

	orderTypeResponse := helper.MapOrderType(&updatedOrderType)

	return *orderTypeResponse, nil
}

func (s *orderTypeService) Delete(id uuid.UUID) error {
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

func (s *orderTypeService) FindById(id uuid.UUID) response.OrderTypeResponse {
	orderTypeData, err := s.repo.FindById(id)
	helper.ErrorPanic(err)

	orderTypeResponse := helper.MapOrderType(&orderTypeData)
	return *orderTypeResponse
}

func (s *orderTypeService) FindWithPagination(pagination request.Pagination) ([]response.OrderTypeResponse, int64, error) {
	orderTypees, total, err := s.repo.FindWithPagination(pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.OrderTypeResponse
	for _, orderType := range orderTypees {
		result = append(result, *helper.MapOrderType(&orderType))
	}

	return result, total, nil
}

func (s *orderTypeService) FindWithPaginationCursor(pagination request.Pagination) ([]response.OrderTypeResponse, string, bool, error) {
	orderTypes, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.OrderTypeResponse
	for _, orderType := range orderTypes {
		result = append(result, *helper.MapOrderType(&orderType))
	}

	return result, nextCursor, hasNext, nil
}
