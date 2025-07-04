package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type TaxService interface {
	Create(req request.TaxCreate) (response.TaxResponse, error)
	Update(id int, req request.TaxUpdate) (response.TaxResponse, error)
	Delete(id int) error
	FindById(roleId int) response.TaxResponse
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.TaxResponse, int64, error)
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

func (s *taxService) Create(req request.TaxCreate) (response.TaxResponse, error) {
	// Validasi input
	if err := s.validate.Struct(req); err != nil {
		return response.TaxResponse{}, err
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	// Buat entity Tax
	tax := entity.Tax{
		BusinessId:   req.BusinessId,
		Name:         req.Name,
		IsPercentage: isPercentageVal,
		Amount:       req.Amount,
		IsGlobal:     req.IsGlobal,
	}

	createdTax, err := s.repo.Create(tax)
	if err != nil {
		return response.TaxResponse{}, err
	}

	// Mapping ke response
	taxResponse := helper.MapTax(&createdTax)

	return *taxResponse, nil
}

func (s *taxService) Update(id int, req request.TaxUpdate) (response.TaxResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TaxResponse{}, err
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	tax := entity.Tax{
		Id:           id,
		Name:         req.Name,
		IsPercentage: isPercentageVal,
		Amount:       req.Amount,
		IsGlobal:     req.IsGlobal,
	}

	updatedTax, err := s.repo.Update(tax)
	if err != nil {
		return response.TaxResponse{}, err
	}

	// Mapping ke response
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
