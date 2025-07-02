package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type BundleService interface {
	CreateBundle(req request.BundleCreate) error
	UpdateBundle(id int, req request.BundleUpdate) error
	FindById(id int) (response.BundleResponse, error)
	Delete(id int) error
	SetIsActive(id int, active bool) error
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.BundleResponse, int64, error) // <- Tambahan
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

func (s *bundleService) CreateBundle(req request.BundleCreate) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	// Upload gambar produk utama
	var imageURL *string
	if req.Image != nil && *req.Image != "" {
		url, err := helper.UploadBase64ToCloudinary(*req.Image, "product")
		if err != nil {
			return fmt.Errorf("gagal upload gambar produk: %w", err)
		}
		imageURL = &url
	}

	bundle := entity.Bundle{
		BusinessId:  req.BusinessId,
		Name:        req.Name,
		Description: req.Description,
		Image:       imageURL,
		BasePrice:   req.BasePrice,
		Stock:       req.Stock,
		TaxId:       req.TaxId,
		IsAvailable: true,
		IsActive:    true,
	}

	if err := s.BundleRepository.InsertBundle(&bundle); err != nil {
		return err
	}

	var items []entity.BundleItem
	for _, item := range req.Items {
		items = append(items, entity.BundleItem{
			BundleId:  bundle.Id,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	return s.BundleRepository.InsertItemsByBundleId(bundle.Id, items)
}

func (s *bundleService) UpdateBundle(id int, req request.BundleUpdate) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	bundle, err := s.BundleRepository.FindById(id)
	if err != nil {
		return err
	}

	bundle.BusinessId = req.BusinessId
	bundle.Name = req.Name
	bundle.Description = req.Description
	bundle.BasePrice = req.BasePrice
	bundle.Stock = req.Stock
	bundle.IsAvailable = req.IsAvailable
	bundle.IsActive = req.IsActive
	bundle.TaxId = req.TaxId

	if err := s.BundleRepository.UpdateBundle(&bundle); err != nil {
		return err
	}

	if err := s.BundleRepository.DeleteItemsByBundleId(bundle.Id); err != nil {
		return err
	}

	var items []entity.BundleItem
	for _, item := range req.Items {
		items = append(items, entity.BundleItem{
			BundleId:  bundle.Id,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	return s.BundleRepository.InsertItemsByBundleId(bundle.Id, items)
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
