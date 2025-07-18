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

type BrandService interface {
	Create(req request.BrandRequest) (response.BrandResponse, error)
	Update(id int, req request.BrandRequest) (response.BrandResponse, error)
	Delete(id int) error
	FindById(roleId int) response.BrandResponse
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.BrandResponse, int64, error)
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
	// Validasi input
	if err := s.validate.Struct(req); err != nil {
		return response.BrandResponse{}, err
	}

	// Buat entity Brand
	brand := entity.Brand{
		BusinessId: req.BusinessId,
		Name:       strings.ToLower(req.Name),
	}

	createdBrand, err := s.repo.Create(brand)
	if err != nil {
		return response.BrandResponse{}, err
	}

	// Mapping ke response
	brandResponse := helper.MapBrand(&createdBrand)

	return *brandResponse, nil
}

func (s *brandService) Update(id int, req request.BrandRequest) (response.BrandResponse, error) {
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

	// Mapping ke response
	brandResponse := helper.MapBrand(&updatedBrand)

	return *brandResponse, nil
}

func (s *brandService) Delete(id int) error {
	brand, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(brand)
}

func (s *brandService) FindById(brandId int) response.BrandResponse {
	brandData, err := s.repo.FindById(brandId)
	helper.ErrorPanic(err)

	brandResponse := helper.MapBrand(&brandData)
	return *brandResponse
}

func (s *brandService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.BrandResponse, int64, error) {
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
