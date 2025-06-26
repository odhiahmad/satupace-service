package service

import (
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

	bundle := entity.Bundle{
		BusinessId:  req.BusinessId,
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		BasePrice:   req.BasePrice,
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
	bundle.Image = req.Image
	bundle.BasePrice = req.BasePrice
	bundle.IsAvailable = req.IsAvailable
	bundle.IsActive = req.IsActive

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
	return mapBundleToResponse(bundle), nil
}

func (s *bundleService) Delete(id int) error {
	return s.BundleRepository.Delete(id)
}

func (s *bundleService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.BundleResponse, int64, error) {
	bundles, total, err := s.BundleRepository.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.BundleResponse
	for _, bundleItem := range bundles {
		result = append(result, mapBundleToResponse(bundleItem))
	}

	return result, total, nil
}

func mapBundleToResponse(p entity.Bundle) response.BundleResponse {
	var items []response.BundleItemResponse
	for _, i := range p.Items {
		items = append(items, response.BundleItemResponse{
			Id:          i.Id,
			ProductId:   i.ProductId,
			Name:        i.Product.Name,
			Description: i.Product.Description,
			Image:       i.Product.Image,
			BasePrice:   i.Product.BasePrice,
			SKU:         i.Product.SKU,
			Stock:       i.Product.Stock,
			IsAvailable: i.Product.IsAvailable,
			IsActive:    i.Product.IsActive,
			Quantity:    i.Quantity,
		})
	}

	description := helper.StringOrDefault(p.Description, "")
	image := helper.StringOrDefault(p.Image, "")

	return response.BundleResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: description,
		Image:       image,
		BasePrice:   p.BasePrice,
		Stock:       p.Stock,
		IsAvailable: p.IsAvailable,
		IsActive:    p.IsActive,
		Items:       items,
	}
}
