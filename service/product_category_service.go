package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ProductCategoryService interface {
	Create(productCategory request.ProductCategoryCreate) error
	Update(productCategory request.ProductCategoryUpdate) error
	FindById(id int) (response.ProductCategoryResponse, error)
	FindAll() ([]response.ProductCategoryResponse, error)
	FindByBusinessId(businessId int) ([]response.ProductCategoryResponse, error)
	Delete(id int) error
}

type ProductCategoryServiceImpl struct {
	Repo     repository.ProductCategoryRepository
	Validate *validator.Validate
}

func NewProductCategoryService(repo repository.ProductCategoryRepository, validate *validator.Validate) ProductCategoryService {
	return &ProductCategoryServiceImpl{
		Repo:     repo,
		Validate: validate,
	}
}

func (s *ProductCategoryServiceImpl) Create(req request.ProductCategoryCreate) error {
	err := s.Validate.Struct(req)
	if err != nil {
		return err
	}

	category := entity.ProductCategory{
		Name:        req.Name,
		ParentId:    req.ParentId,
		BusinessId:  req.BusinessId,
		IsActive:    true,
		IsAvailable: true,
	}

	return s.Repo.InsertProductCategory(category)
}
func (s *ProductCategoryServiceImpl) Update(req request.ProductCategoryUpdate) error {
	err := s.Validate.Struct(req)
	if err != nil {
		return err
	}

	category, err := s.Repo.FindById(req.Id)
	if err != nil {
		return err
	}

	category.Name = req.Name
	category.ParentId = req.ParentId

	return s.Repo.UpdateProductCategory(category)
}

func (s *ProductCategoryServiceImpl) FindById(id int) (response.ProductCategoryResponse, error) {
	category, err := s.Repo.FindById(id)
	if err != nil {
		return response.ProductCategoryResponse{}, err
	}

	return response.ProductCategoryResponse{
		Id:       category.Id,
		Name:     category.Name,
		ParentId: category.ParentId,
	}, nil
}

func (s *ProductCategoryServiceImpl) FindAll() ([]response.ProductCategoryResponse, error) {
	categories, err := s.Repo.FindAll()
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

func (s *ProductCategoryServiceImpl) FindByBusinessId(businessId int) ([]response.ProductCategoryResponse, error) {
	categories, err := s.Repo.FindByBusinessId(businessId)
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
func (s *ProductCategoryServiceImpl) Delete(id int) error {
	return s.Repo.Delete(id)
}
