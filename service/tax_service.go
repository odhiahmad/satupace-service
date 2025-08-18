package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/helper/mapper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type TaxService interface {
	Create(req request.TaxRequest) (response.TaxResponse, error)
	Update(id uuid.UUID, req request.TaxRequest) (response.TaxResponse, error)
	Delete(id uuid.UUID) error
	FindById(roleId uuid.UUID) response.TaxResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.TaxResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.TaxResponse, string, bool, error)
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
		IsPercentage: &isPercentageVal,
		IsGlobal:     req.IsGlobal,
		Amount:       req.Amount,
	}

	createdTax, err := s.repo.Create(tax)
	if err != nil {
		return response.TaxResponse{}, err
	}

	taxResponse := mapper.MapTax(&createdTax)

	return *taxResponse, nil
}

func (s *taxService) Update(id uuid.UUID, req request.TaxRequest) (response.TaxResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TaxResponse{}, err
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	tax := entity.Tax{
		Id:           id,
		Name:         strings.ToLower(req.Name),
		IsPercentage: &isPercentageVal,
		Amount:       req.Amount,
		IsGlobal:     req.IsGlobal,
		IsActive:     req.IsActive,
	}

	updatedTax, err := s.repo.Update(tax)
	if err != nil {
		return response.TaxResponse{}, err
	}

	taxResponse := mapper.MapTax(&updatedTax)

	return *taxResponse, nil
}

func (s *taxService) Delete(id uuid.UUID) error {
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

func (s *taxService) FindById(taxId uuid.UUID) response.TaxResponse {
	taxData, err := s.repo.FindById(taxId)
	helper.ErrorPanic(err)

	taxResponse := mapper.MapTax(&taxData)
	return *taxResponse
}

func (s *taxService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.TaxResponse, int64, error) {
	taxes, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.TaxResponse
	for _, tax := range taxes {
		result = append(result, *mapper.MapTax(&tax))
	}

	return result, total, nil
}

func (s *taxService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.TaxResponse, string, bool, error) {
	taxes, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.TaxResponse
	for _, tax := range taxes {
		result = append(result, *mapper.MapTax(&tax))
	}

	return result, nextCursor, hasNext, nil
}
