package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type CategoryService interface {
	Create(category request.CategoryRequest) (response.CategoryResponse, error)
	Update(id int, category request.CategoryRequest) (response.CategoryResponse, error)
	FindById(brandId int) response.CategoryResponse
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.CategoryResponse, int64, error)
	Delete(id int) error
}

type categoryService struct {
	repo     repository.CategoryRepository
	Validate *validator.Validate
}

func NewCategoryService(repo repository.CategoryRepository, validate *validator.Validate) CategoryService {
	return &categoryService{
		repo:     repo,
		Validate: validate,
	}
}

func (s *categoryService) Create(req request.CategoryRequest) (response.CategoryResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	category := entity.Category{
		Name:       req.Name,
		ParentId:   req.ParentId,
		BusinessId: req.BusinessId,
		IsActive:   true,
	}

	createdCategory, err := s.repo.InsertCategory(category)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	categoryResponse := helper.MapCategory(&createdCategory)
	return *categoryResponse, nil
}

func (s *categoryService) Update(id int, req request.CategoryRequest) (response.CategoryResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	// Ambil data lama
	category, err := s.repo.FindById(id)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	// Update field yang boleh diubah
	category.Name = req.Name
	category.BusinessId = req.BusinessId
	category.ParentId = req.ParentId

	updatedCategory, err := s.repo.UpdateCategory(category)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	categoryResponse := helper.MapCategory(&updatedCategory)
	return *categoryResponse, nil
}

func (s *categoryService) FindById(brandId int) response.CategoryResponse {
	categories, err := s.repo.FindById(brandId)
	helper.ErrorPanic(err)

	category := helper.MapCategory(&categories)
	return *category
}

func (s *categoryService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.CategoryResponse, int64, error) {
	categories, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.CategoryResponse
	for _, category := range categories {
		result = append(result, *helper.MapCategory(&category))
	}

	return result, total, nil
}

func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
