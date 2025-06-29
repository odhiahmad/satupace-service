package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type BusinessService interface {
	Create(req request.BusinessCreate) (response.BusinessResponse, error)
	Update(req request.BusinessUpdate) (response.BusinessResponse, error)
	Delete(id int) error
	FindById(id int) (response.BusinessResponse, error)
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
		Name:           req.Name,
		OwnerName:      req.OwnerName,
		BusinessTypeId: req.BusinessTypeId,
		Image:          req.Image,
		IsActive:       req.IsActive,
	}

	created, err := s.repo.Create(business)
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return MapToBusinessResponse(created), nil
}

func (s *businessService) Update(req request.BusinessUpdate) (response.BusinessResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.BusinessResponse{}, err
	}

	business := entity.Business{
		Id:             req.Id,
		Name:           req.Name,
		OwnerName:      req.OwnerName,
		BusinessTypeId: req.BusinessTypeId,
		Image:          req.Image,
		IsActive:       req.IsActive,
	}

	updated, err := s.repo.Update(business)
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return MapToBusinessResponse(updated), nil
}

func (s *businessService) Delete(id int) error {
	business, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(business)
}

func (s *businessService) FindById(id int) (response.BusinessResponse, error) {
	business, err := s.repo.FindById(id)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	return MapToBusinessResponse(business), nil
}

func (s *businessService) FindWithPagination(pagination request.Pagination) ([]response.BusinessResponse, int64, error) {
	businesses, total, err := s.repo.FindWithPagination(pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []response.BusinessResponse
	for _, b := range businesses {
		responses = append(responses, MapToBusinessResponse(b))
	}

	return responses, total, nil
}

func MapToBusinessResponse(b entity.Business) response.BusinessResponse {
	var branches []response.BusinessBranchResponse
	for _, branch := range b.Branches {
		branches = append(branches, response.BusinessBranchResponse{
			Id:          branch.Id,
			BusinessId:  branch.BusinessId,
			PhoneNumber: helper.StringValue(branch.PhoneNumber),
			Rating:      helper.StringValue(branch.Rating),
			Provinsi:    helper.StringValue(branch.Provinsi),
			Kota:        helper.StringValue(branch.Kota),
			Kecamatan:   helper.StringValue(branch.Kecamatan),
			PostalCode:  helper.StringValue(branch.PostalCode),
			IsMain:      branch.IsMain,
			IsActive:    branch.IsActive,
		})
	}

	return response.BusinessResponse{
		Id:             b.Id,
		Name:           b.Name,
		OwnerName:      b.OwnerName,
		BusinessTypeId: b.BusinessTypeId,
		Image:          b.Image,
		IsActive:       b.IsActive,
		Branches:       branches,
	}
}
