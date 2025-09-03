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
)

type OrderTypeService interface {
	Create(req request.OrderTypeRequest) (response.OrderTypeResponse, error)
	Update(id int, req request.OrderTypeRequest) (response.OrderTypeResponse, error)
	Delete(id int) error
	FindById(id int) response.OrderTypeResponse
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

	orderTypeResponse := mapper.MapOrderType(&createdOrderType)

	return *orderTypeResponse, nil
}

func (s *orderTypeService) Update(id int, req request.OrderTypeRequest) (response.OrderTypeResponse, error) {
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

	orderTypeResponse := mapper.MapOrderType(&updatedOrderType)

	return *orderTypeResponse, nil
}

func (s *orderTypeService) Delete(id int) error {
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

func (s *orderTypeService) FindById(id int) response.OrderTypeResponse {
	orderTypeData, err := s.repo.FindById(id)
	helper.ErrorPanic(err)

	orderTypeResponse := mapper.MapOrderType(&orderTypeData)
	return *orderTypeResponse
}

func (s *orderTypeService) FindWithPagination(pagination request.Pagination) ([]response.OrderTypeResponse, int64, error) {
	orderTypees, total, err := s.repo.FindWithPagination(pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.OrderTypeResponse
	for _, orderType := range orderTypees {
		result = append(result, *mapper.MapOrderType(&orderType))
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
		result = append(result, *mapper.MapOrderType(&orderType))
	}

	return result, nextCursor, hasNext, nil
}
