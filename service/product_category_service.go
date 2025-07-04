package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ProductCategoryService interface {
	Create(productCategory request.ProductCategoryCreate) (response.ProductCategoryResponse, error)
	Update(productCategory request.ProductCategoryUpdate) (response.ProductCategoryResponse, error)
	FindById(id int) (response.ProductCategoryResponse, error)
	FindAll() ([]response.ProductCategoryResponse, error)
	FindByBusinessId(businessId int) ([]response.ProductCategoryResponse, error)
	Delete(id int) error
}

type productCategoryService struct {
	repo     repository.ProductCategoryRepository
	Validate *validator.Validate
}

func NewProductCategoryService(repo repository.ProductCategoryRepository, validate *validator.Validate) ProductCategoryService {
	return &productCategoryService{
		repo:     repo,
		Validate: validate,
	}
}

func (s *productCategoryService) Create(req request.ProductCategoryCreate) (response.ProductCategoryResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return response.ProductCategoryResponse{}, err
	}

	category := entity.ProductCategory{
		Name:       req.Name,
		ParentId:   req.ParentId,
		BusinessId: req.BusinessId,
		IsActive:   true,
	}

	createdCategory, err := s.repo.InsertProductCategory(category)
	if err != nil {
		return response.ProductCategoryResponse{}, err
	}

	categoryResponse := helper.MapProductCategory(&createdCategory)
	return *categoryResponse, nil
}

func (s *productCategoryService) Update(req request.ProductCategoryUpdate) (response.ProductCategoryResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return response.ProductCategoryResponse{}, err
	}

	// Ambil data lama
	category, err := s.repo.FindById(req.Id)
	if err != nil {
		return response.ProductCategoryResponse{}, err
	}

	// Update field yang boleh diubah
	category.Name = req.Name
	category.ParentId = req.ParentId

	updatedCategory, err := s.repo.UpdateProductCategory(category)
	if err != nil {
		return response.ProductCategoryResponse{}, err
	}

	categoryResponse := helper.MapProductCategory(&updatedCategory)
	return *categoryResponse, nil
}

func (s *productCategoryService) FindById(id int) (response.ProductCategoryResponse, error) {
	category, err := s.repo.FindById(id)
	if err != nil {
		return response.ProductCategoryResponse{}, err
	}

	return response.ProductCategoryResponse{
		Id:       category.Id,
		Name:     category.Name,
		ParentId: category.ParentId,
	}, nil
}

func (s *productCategoryService) FindAll() ([]response.ProductCategoryResponse, error) {
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []response.ProductCategoryResponse
	for _, c := range categories {
		result = append(result, response.ProductCategoryResponse{
			Id:       c.Id,
			Name:     c.Name,
			ParentId: c.ParentId,
		})
	}

	return result, nil
}

func (s *productCategoryService) FindByBusinessId(businessId int) ([]response.ProductCategoryResponse, error) {
	categories, err := s.repo.FindByBusinessId(businessId)
	if err != nil {
		return nil, err
	}

	var result []response.ProductCategoryResponse
	for _, c := range categories {
		result = append(result, response.ProductCategoryResponse{
			Id:       c.Id,
			Name:     c.Name,
			ParentId: c.ParentId,
		})
	}

	return result, nil
}
func (s *productCategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
