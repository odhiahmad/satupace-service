package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type BundleService interface {
	CreateBundle(req request.BundleCreate)
	UpdateBundle(req request.BundleUpdate)
	FindById(bundleId int) response.BundleResponse
	FindAll() []response.BundleResponse
	Delete(bundleId int)
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

func (s *BundleServiceImpl) CreateBundle(req request.BundleCreate) {
	err := s.Validate.Struct(req)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	entityBundle := entity.Bundle{
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

	err = entityBundle.Prepare()
	helper.ErrorPanic(err)

	s.BundleRepository.InsertBundle(entityBundle)
	saved := entityBundle

	var items []entity.BundleItem
	for _, item := range req.Items {
		items = append(items, entity.BundleItem{
			BundleId:  saved.Id,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	s.BundleRepository.InsertItemsByBundleId(saved.Id, items)
}

func (s *BundleServiceImpl) UpdateBundle(req request.BundleUpdate) {
	err := s.Validate.Struct(req)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	existing, err := s.BundleRepository.FindById(req.Id)
	helper.ErrorPanic(err)

	existing.BusinessId = req.BusinessId
	existing.Name = req.Name
	existing.Description = req.Description
	existing.Image = req.Image
	existing.BasePrice = req.BasePrice
	existing.FinalPrice = req.FinalPrice
	existing.Discount = req.Discount
	existing.Promo = req.Promo
	existing.IsAvailable = req.IsAvailable
	existing.IsActive = req.IsActive

	s.BundleRepository.UpdateBundle(existing)
	s.BundleRepository.DeleteItemsByBundleId(existing.Id)

	var items []entity.BundleItem
	for _, item := range req.Items {
		items = append(items, entity.BundleItem{
			BundleId:  existing.Id,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	s.BundleRepository.InsertItemsByBundleId(existing.Id, items)
}

func (s *BundleServiceImpl) FindById(bundleId int) response.BundleResponse {
	data, err := s.BundleRepository.FindById(bundleId)
	helper.ErrorPanic(err)
	return mapBundleToResponse(data)
}

func (s *BundleServiceImpl) FindAll() []response.BundleResponse {
	products := s.BundleRepository.FindAll()
	var responses []response.BundleResponse
	for _, p := range products {
		responses = append(responses, mapBundleToResponse(p))
	}
	return responses
}

func (s *BundleServiceImpl) Delete(bundleId int) {
	s.BundleRepository.Delete(bundleId)
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
