package service

import (
	"strings"

	"loka-kasir/data/request"
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper/mapper"
	"loka-kasir/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BusinessService interface {
	Create(req request.BusinessCreate) (response.BusinessResponse, error)
	Update(req request.BusinessUpdate) (response.BusinessResponse, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (response.BusinessResponse, error)
	FindWithPagination(pagination request.Pagination) ([]response.BusinessResponse, int64, error)
}

type businessService struct {
	repo     repository.BusinessRepository
	validate *validator.Validate
}

func NewBusinessService(repo repository.BusinessRepository, validate *validator.Validate) BusinessService {
	return &businessService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *businessService) Create(req request.BusinessCreate) (response.BusinessResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BusinessResponse{}, err
	}

	business := entity.Business{
		Name:           strings.ToLower(req.Name),
		OwnerName:      req.OwnerName,
		BusinessTypeId: &req.BusinessTypeId,
		Image:          req.Image,
		IsActive:       req.IsActive,
	}

	created, err := s.repo.Create(business)
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return *mapper.MapBusiness(&created), nil
}
func (s *businessService) Update(req request.BusinessUpdate) (response.BusinessResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BusinessResponse{}, err
	}

	business := entity.Business{
		Id:             req.Id,
		Name:           strings.ToLower(req.Name),
		OwnerName:      req.OwnerName,
		BusinessTypeId: &req.BusinessTypeId,
		ProvinceID:     &req.ProvinceID,
		CityID:         &req.CityID,
		DistrictID:     &req.DistrictID,
		VillageID:      &req.VillageID,
		IsActive:       req.IsActive,
	}

	updated, err := s.repo.Update(business)
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return *mapper.MapBusiness(&updated), nil
}

func (s *businessService) Delete(id uuid.UUID) error {
	business, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(business)
}

func (s *businessService) FindById(id uuid.UUID) (response.BusinessResponse, error) {
	business, err := s.repo.FindById(id)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	return *mapper.MapBusiness(&business), nil
}

func (s *businessService) FindWithPagination(pagination request.Pagination) ([]response.BusinessResponse, int64, error) {
	businesses, total, err := s.repo.FindWithPagination(pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []response.BusinessResponse
	for _, b := range businesses {
		responses = append(responses, *mapper.MapBusiness(&b))
	}

	return responses, total, nil
}
