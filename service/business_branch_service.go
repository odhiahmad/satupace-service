package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type BusinessBranchService interface {
	Create(req request.BusinessBranchCreate) (entity.BusinessBranch, error)
	Update(req request.BusinessBranchUpdate) (entity.BusinessBranch, error)
	Delete(id int) error
	FindById(id int) (entity.BusinessBranch, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.BusinessBranch, int64, error)
}

type businessBranchService struct {
	repo     repository.BusinessBranchRepository
	validate *validator.Validate
}

func NewBusinessBranchService(repo repository.BusinessBranchRepository, validate *validator.Validate) BusinessBranchService {
	return &businessBranchService{
		repo:     repo,
		validate: validator.New(), // ensure always uses a fresh validator
	}
}

func (s *businessBranchService) Create(req request.BusinessBranchCreate) (entity.BusinessBranch, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.BusinessBranch{}, err
	}

	businessBranch := entity.BusinessBranch{
		BusinessId:  req.BusinessId,
		PhoneNumber: req.PhoneNumber,
		Rating:      req.Rating,
		Provinsi:    req.Provinsi,
		Kota:        req.Kota,
		Kecamatan:   req.Kecamatan,
		PostalCode:  req.PostalCode,
		Phone:       req.Phone,
		IsMain:      req.IsMain,
		IsActive:    req.IsActive,
	}

	return s.repo.Create(businessBranch)
}

func (s *businessBranchService) Update(req request.BusinessBranchUpdate) (entity.BusinessBranch, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.BusinessBranch{}, err
	}

	businessBranch := entity.BusinessBranch{
		Id:          req.Id,
		PhoneNumber: req.PhoneNumber,
		Rating:      req.Rating,
		Provinsi:    req.Provinsi,
		Kota:        req.Kota,
		Kecamatan:   req.Kecamatan,
		PostalCode:  req.PostalCode,
		Phone:       req.Phone,
		IsMain:      req.IsMain,
		IsActive:    req.IsActive,
	}

	return s.repo.Update(businessBranch)
}

func (s *businessBranchService) Delete(id int) error {
	businessBranch, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(businessBranch)
}

func (s *businessBranchService) FindById(id int) (entity.BusinessBranch, error) {
	return s.repo.FindById(id)
}

func (s *businessBranchService) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.BusinessBranch, int64, error) {
	return s.repo.FindWithPagination(businessId, pagination)
}

func MapToBusinessBranchResponse(branch entity.BusinessBranch) response.BusinessBranchResponse {
	return response.BusinessBranchResponse{
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
	}
}
