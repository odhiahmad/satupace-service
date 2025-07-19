package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type TaxService interface {
	Create(req request.TaxRequest) (response.TaxResponse, error)
	Update(id int, req request.TaxRequest) (response.TaxResponse, error)
	Delete(id int) error
	FindById(roleId int) response.TaxResponse
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.TaxResponse, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.TaxResponse, string, bool, error)
}

type taxService struct {
	repo     repository.TaxRepository
	validate *validator.Validate
}

func NewTaxService(repo repository.TaxRepository, validate *validator.Validate) TaxService {
	return &taxService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *taxService) Create(req request.TaxRequest) (response.TaxResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TaxResponse{}, err
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	tax := entity.Tax{
		BusinessId:   req.BusinessId,
		Name:         strings.ToLower(req.Name),
		IsPercentage: isPercentageVal,
		Amount:       req.Amount,
	}

	createdTax, err := s.repo.Create(tax)
	if err != nil {
		return response.TaxResponse{}, err
	}

	taxResponse := helper.MapTax(&createdTax)

	return *taxResponse, nil
}

func (s *taxService) Update(id int, req request.TaxRequest) (response.TaxResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TaxResponse{}, err
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	tax := entity.Tax{
		Id:           id,
		Name:         strings.ToLower(req.Name),
		IsPercentage: isPercentageVal,
		Amount:       req.Amount,
	}

	updatedTax, err := s.repo.Update(tax)
	if err != nil {
		return response.TaxResponse{}, err
	}

	taxResponse := helper.MapTax(&updatedTax)

	return *taxResponse, nil
}

func (s *taxService) Delete(id int) error {
	tax, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(tax)
}

func (s *taxService) FindById(taxId int) response.TaxResponse {
	taxData, err := s.repo.FindById(taxId)
	helper.ErrorPanic(err)

	taxResponse := helper.MapTax(&taxData)
	return *taxResponse
}

func (s *taxService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.TaxResponse, int64, error) {
	taxes, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.TaxResponse
	for _, tax := range taxes {
		result = append(result, *helper.MapTax(&tax))
	}

	return result, total, nil
}

func (s *taxService) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.TaxResponse, string, bool, error) {
	taxes, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.TaxResponse
	for _, tax := range taxes {
		result = append(result, *helper.MapTax(&tax))
	}

	return result, nextCursor, hasNext, nil
}
