package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type BundleService interface {
	CreateBundle(req request.BundleCreate) error
	UpdateBundle(id int, req request.BundleUpdate) error
	FindById(id int) (response.BundleResponse, error)
	FindAll() ([]response.BundleResponse, error)
	FindByBusinessId(businessId int) ([]response.BundleResponse, error)
	Delete(id int) error
}

type BundleServiceImpl struct {
	BundleRepository repository.BundleRepository
	Validate         *validator.Validate
}

func NewBundleService(repo repository.BundleRepository, validate *validator.Validate) BundleService {
	return &BundleServiceImpl{
		BundleRepository: repo,
		Validate:         validate,
	}
}

func (s *BundleServiceImpl) CreateBundle(req request.BundleCreate) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	bundle := entity.Bundle{
		BusinessId:  req.BusinessId,
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		BasePrice:   req.BasePrice,
		FinalPrice:  req.FinalPrice,
		Discount:    req.Discount,
		Promo:       req.Promo,
		IsAvailable: true,
		IsActive:    true,
	}

	if err := bundle.Prepare(); err != nil {
		return err
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

func (s *BundleServiceImpl) UpdateBundle(id int, req request.BundleUpdate) error {
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
	bundle.FinalPrice = req.FinalPrice
	bundle.Discount = req.Discount
	bundle.Promo = req.Promo
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

func (s *BundleServiceImpl) FindById(id int) (response.BundleResponse, error) {
	bundle, err := s.BundleRepository.FindById(id)
	if err != nil {
		return response.BundleResponse{}, err
	}
	return mapBundleToResponse(bundle), nil
}

func (s *BundleServiceImpl) FindAll() ([]response.BundleResponse, error) {
	bundles, err := s.BundleRepository.FindAll()
	if err != nil {
		return nil, err
	}
	var responses []response.BundleResponse
	for _, b := range bundles {
		responses = append(responses, mapBundleToResponse(b))
	}
	return responses, nil
}

func (s *BundleServiceImpl) FindByBusinessId(businessId int) ([]response.BundleResponse, error) {
	bundles, err := s.BundleRepository.FindByBusinessId(businessId)
	if err != nil {
		return nil, err
	}

	var responses []response.BundleResponse
	for _, b := range bundles {
		responses = append(responses, mapBundleToResponse(b))
	}

	return responses, nil
}

func (s *BundleServiceImpl) Delete(id int) error {
	return s.BundleRepository.Delete(id)
}

func mapBundleToResponse(p entity.Bundle) response.BundleResponse {
	var items []response.BundleItemResponse
	for _, i := range p.Items {
		items = append(items, response.BundleItemResponse{
			Id:        i.Id,
			ProductId: i.ProductId,
			Product:   i.Product.Name,
			Quantity:  i.Quantity,
		})
	}

	return response.BundleResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Image:       p.Image,
		BasePrice:   p.BasePrice,
		FinalPrice:  p.FinalPrice,
		Discount:    p.Discount,
		Promo:       p.Promo,
		Stock:       p.Stock,
		IsAvailable: p.IsAvailable,
		IsActive:    p.IsActive,
		Items:       items,
	}
}
