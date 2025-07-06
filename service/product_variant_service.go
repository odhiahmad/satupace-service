package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ProductVariantService interface {
	Create(req request.ProductVariantRequest, productId int) (*entity.ProductVariant, error)
	Update(id int, req request.ProductVariantRequest) (*entity.ProductVariant, error)
	Delete(id int) error
	DeleteByProductId(productId int) error
	FindById(id int) (*entity.ProductVariant, error)
	FindByProductId(productId int) ([]entity.ProductVariant, error)
	SetActive(id int, isActive bool) error
	SetAvailable(id int, isAvailable bool) error
}

type productVariantService struct {
	repo        repository.ProductVariantRepository
	productRepo repository.ProductRepository // Tambahkan ini
	validate    *validator.Validate
}

func NewProductVariantService(repo repository.ProductVariantRepository, productRepo repository.ProductRepository, validate *validator.Validate) ProductVariantService {
	return &productVariantService{
		repo:        repo,
		productRepo: productRepo,
		validate:    validate,
	}
}

func (s *productVariantService) Create(req request.ProductVariantRequest, productId int) (*entity.ProductVariant, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindById(productId)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	sku := req.SKU
	trackStock := req.Stock == 0

	if sku == "" {
		sku = helper.GenerateSKU(req.Name)
	}

	variant := entity.ProductVariant{
		ProductId:  productId,
		BusinessId: req.BusinessId,
		Name:       req.Name,
		BasePrice:  req.BasePrice,
		SKU:        sku,
		Stock:      req.Stock,
		TrackStock: trackStock,
	}

	err = s.repo.Create(&variant)
	if err != nil {
		return nil, err
	}

	if !product.HasVariant {
		_ = s.productRepo.SetHasVariant(productId)
	}

	return &variant, nil
}

func (s *productVariantService) Update(id int, req request.ProductVariantRequest) (*entity.ProductVariant, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	sku := req.SKU

	if sku == "" {
		sku = helper.GenerateSKU(req.Name)
	}

	trackStock := req.Stock == 0

	// Update field
	existing.Name = req.Name
	existing.BasePrice = req.BasePrice
	existing.SKU = sku
	existing.Stock = req.Stock
	existing.TrackStock = trackStock

	err = s.repo.Update(&existing)
	if err != nil {
		return nil, err
	}

	// TODO: Update relasi many2many ke promos jika diperlukan

	return &existing, err
}

func (s *productVariantService) Delete(id int) error {
	// 1. Ambil variant sebelum dihapus untuk dapatkan ProductId
	variant, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	// 2. Hapus variant
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 3. Hitung sisa variant pada produk tersebut
	count, err := s.repo.CountByProductId(variant.ProductId)
	if err != nil {
		return err
	}

	// 4. Jika tidak ada variant lagi, reset status produk
	if count == 0 {
		if err := s.productRepo.ResetVariantStateToFalse(variant.ProductId); err != nil {
			return err
		}
	}

	return nil
}

func (s *productVariantService) DeleteByProductId(productId int) error {
	return s.repo.DeleteByProductId(productId)
}

func (s *productVariantService) FindById(id int) (*entity.ProductVariant, error) {
	variant, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (s *productVariantService) FindByProductId(productId int) ([]entity.ProductVariant, error) {
	return s.repo.FindByProductId(productId)
}

func (s *productVariantService) SetActive(id int, isActive bool) error {
	return s.repo.SetActive(id, isActive)
}

func (s *productVariantService) SetAvailable(id int, isAvailable bool) error {
	return s.repo.SetAvailable(id, isAvailable)
}
