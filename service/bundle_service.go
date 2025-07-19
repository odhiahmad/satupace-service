package service

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type BundleService interface {
	CreateBundle(req request.BundleRequest) (response.BundleResponse, error)
	UpdateBundle(id int, req request.BundleRequest) (response.BundleResponse, error)
	FindById(id int) (response.BundleResponse, error)
	Delete(id int) error
	SetIsActive(id int, active bool) error
	SetIsAvailable(id int, isAvailable bool) error
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.BundleResponse, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.BundleResponse, string, bool, error)
}

type bundleService struct {
	BundleRepository repository.BundleRepository
	Validate         *validator.Validate
}

func NewBundleService(repo repository.BundleRepository, validate *validator.Validate) BundleService {
	return &bundleService{
		BundleRepository: repo,
		Validate:         validate,
	}
}

func (s *bundleService) CreateBundle(req request.BundleRequest) (response.BundleResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return response.BundleResponse{}, err
	}

	// Upload gambar produk utama
	var imageURL *string
	if req.Image != nil && *req.Image != "" {
		url, err := helper.UploadBase64ToCloudinary(*req.Image, "product")
		if err != nil {
			return response.BundleResponse{}, fmt.Errorf("gagal upload gambar produk: %w", err)
		}
		imageURL = &url
	}

	bundle := entity.Bundle{
		BusinessId:  req.BusinessId,
		Name:        strings.ToLower(req.Name),
		Description: req.Description,
		Image:       imageURL,
		BasePrice:   req.BasePrice,
		Stock:       req.Stock,
		TaxId:       req.TaxId,
		IsAvailable: true,
		IsActive:    true,
	}

	createdBundle, err := s.BundleRepository.InsertBundle(&bundle)
	if err != nil {
		return response.BundleResponse{}, err
	}

	var items []entity.BundleItem
	for _, item := range req.Items {
		items = append(items, entity.BundleItem{
			BundleId:  bundle.Id,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	err = s.BundleRepository.InsertItemsByBundleId(bundle.Id, items)
	if err != nil {
		return response.BundleResponse{}, err
	}

	bundleResponse := helper.MapBundleToResponse(createdBundle)

	return bundleResponse, nil
}

func (s *bundleService) UpdateBundle(id int, req request.BundleRequest) (response.BundleResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return response.BundleResponse{}, err
	}

	bundle, err := s.BundleRepository.FindById(id)
	if err != nil {
		return response.BundleResponse{}, err
	}

	bundle.BusinessId = req.BusinessId
	bundle.Name = strings.ToLower(req.Name)
	bundle.Description = req.Description
	bundle.BasePrice = req.BasePrice
	bundle.Stock = req.Stock
	bundle.TaxId = req.TaxId

	updatedBundle, err := s.BundleRepository.UpdateBundle(&bundle)
	if err != nil {
		return response.BundleResponse{}, err
	}

	if err := s.BundleRepository.DeleteItemsByBundleId(bundle.Id); err != nil {
		return response.BundleResponse{}, err
	}

	var items []entity.BundleItem
	for _, item := range req.Items {
		items = append(items, entity.BundleItem{
			BundleId:  bundle.Id,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	err = s.BundleRepository.InsertItemsByBundleId(bundle.Id, items)
	if err != nil {
		return response.BundleResponse{}, err
	}

	bundleResponse := helper.MapBundleToResponse(updatedBundle)

	return bundleResponse, nil
}

func (s *bundleService) FindById(id int) (response.BundleResponse, error) {
	bundle, err := s.BundleRepository.FindById(id)
	if err != nil {
		return response.BundleResponse{}, err
	}
	return helper.MapBundleToResponse(bundle), nil
}

func (s *bundleService) Delete(id int) error {
	return s.BundleRepository.Delete(id)
}

func (s *bundleService) SetIsActive(id int, active bool) error {
	return s.BundleRepository.SetIsActive(id, active)
}

func (s *bundleService) SetIsAvailable(id int, active bool) error {
	return s.BundleRepository.SetIsAvailable(id, active)
}

func (s *bundleService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.BundleResponse, int64, error) {
	bundles, total, err := s.BundleRepository.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.BundleResponse
	for _, bundleItem := range bundles {
		result = append(result, helper.MapBundleToResponse(bundleItem))
	}

	return result, total, nil
}

func (s *bundleService) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.BundleResponse, string, bool, error) {
	bundles, nextCursor, hasNext, err := s.BundleRepository.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.BundleResponse
	for _, bundleItem := range bundles {
		result = append(result, helper.MapBundleToResponse(bundleItem))
	}

	return result, nextCursor, hasNext, nil
}
