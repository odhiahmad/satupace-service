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

type BrandService interface {
	Create(req request.BrandRequest) (response.BrandResponse, error)
	Update(id uuid.UUID, req request.BrandRequest) (response.BrandResponse, error)
	Delete(id uuid.UUID) error
	FindById(roleId uuid.UUID) response.BrandResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.BrandResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.BrandResponse, string, bool, error)
}

type brandService struct {
	repo     repository.BrandRepository
	validate *validator.Validate
}

func NewBrandService(repo repository.BrandRepository, validate *validator.Validate) BrandService {
	return &brandService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *brandService) Create(req request.BrandRequest) (response.BrandResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BrandResponse{}, err
	}

	brand := entity.Brand{
		BusinessId: req.BusinessId,
		Name:       strings.ToLower(req.Name),
	}

	createdBrand, err := s.repo.Create(brand)
	if err != nil {
		return response.BrandResponse{}, err
	}

	brandResponse := helper.MapBrand(&createdBrand)

	return *brandResponse, nil
}

func (s *brandService) Update(id uuid.UUID, req request.BrandRequest) (response.BrandResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BrandResponse{}, err
	}

	brand := entity.Brand{
		Id:   id,
		Name: strings.ToLower(req.Name),
	}

	updatedBrand, err := s.repo.Update(brand)
	if err != nil {
		return response.BrandResponse{}, err
	}

	brandResponse := helper.MapBrand(&updatedBrand)

	return *brandResponse, nil
}

func (s *brandService) Delete(id uuid.UUID) error {
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

func (s *brandService) FindById(brandId uuid.UUID) response.BrandResponse {
	brandData, err := s.repo.FindById(brandId)
	helper.ErrorPanic(err)

	brandResponse := helper.MapBrand(&brandData)
	return *brandResponse
}

func (s *brandService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.BrandResponse, int64, error) {
	brandes, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.BrandResponse
	for _, brand := range brandes {
		result = append(result, *helper.MapBrand(&brand))
	}

	return result, total, nil
}

func (s *brandService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.BrandResponse, string, bool, error) {
	brands, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.BrandResponse
	for _, brand := range brands {
		result = append(result, *helper.MapBrand(&brand))
	}

	return result, nextCursor, hasNext, nil
}
