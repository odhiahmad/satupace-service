package service

import (
	"fmt"
	"strings"

	"loka-kasir/data/request"
	"loka-kasir/entity"
	"loka-kasir/helper"
	"loka-kasir/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProductVariantService interface {
	Create(req request.ProductVariantRequest, productId uuid.UUID) (*entity.ProductVariant, error)
	Update(id uuid.UUID, req request.ProductVariantRequest) (*entity.ProductVariant, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (*entity.ProductVariant, error)
	FindByProductId(productId uuid.UUID) ([]entity.ProductVariant, error)
	SetActive(id uuid.UUID, isActive bool) error
	SetAvailable(id uuid.UUID, isAvailable bool) error
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

func (s *productVariantService) Create(req request.ProductVariantRequest, productId uuid.UUID) (*entity.ProductVariant, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindById(productId)
	if err != nil {
		return nil, fmt.Errorf("product tidak ditemukan: %w", err)
	}

	var sku string
	if req.SKU == nil || *req.SKU == "" {
		generatedSKU := helper.GenerateSKU(strings.ToLower(req.Name))
		sku = generatedSKU
	} else {
		sku = *req.SKU
	}

	exist, err := s.repo.IsSKUExist(sku, req.BusinessId)
	if err != nil {
		return nil, fmt.Errorf("gagal cek SKU: %w", err)
	}
	if exist {
		return nil, fmt.Errorf("SKU varian sudah digunakan: %s", sku)
	}

	variant := entity.ProductVariant{
		ProductId:        &productId,
		BusinessId:       req.BusinessId,
		Name:             strings.ToLower(req.Name),
		Description:      helper.LowerStringPtr(req.Description),
		BasePrice:        req.BasePrice,
		SellPrice:        req.SellPrice,
		SKU:              &sku,
		Stock:            req.Stock,
		TrackStock:       req.TrackStock,
		IgnoreStockCheck: req.IgnoreStockCheck,
		IsAvailable:      req.IsAvailable,
		IsActive:         req.IsActive,
		MinimumSales:     req.MinimumSales,
	}

	if err := s.repo.Create(&variant); err != nil {
		return nil, fmt.Errorf("gagal menyimpan varian: %w", err)
	}

	if !product.HasVariant {
		if err := s.productRepo.SetHasVariant(productId); err != nil {
			return nil, fmt.Errorf("gagal mengupdate hasVariant produk: %w", err)
		}
	}

	return &variant, nil
}

func (s *productVariantService) Update(id uuid.UUID, req request.ProductVariantRequest) (*entity.ProductVariant, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("variant tidak ditemukan: %w", err)
	}

	sku := req.SKU
	if sku == nil || *sku == "" {
		s := helper.GenerateSKU(strings.ToLower(req.Name))
		sku = &s
	}

	exist, err := s.repo.IsSKUExistExcept(*sku, req.BusinessId, id)
	if err != nil {
		return nil, fmt.Errorf("gagal cek SKU: %w", err)
	}
	if exist {
		return nil, fmt.Errorf("SKU sudah digunakan oleh variant lain")
	}

	existing.Name = strings.ToLower(req.Name)
	existing.BasePrice = req.BasePrice
	existing.SellPrice = req.SellPrice
	existing.SKU = sku
	existing.Stock = req.Stock
	existing.TrackStock = req.TrackStock

	if err := s.repo.Update(&existing); err != nil {
		return nil, fmt.Errorf("gagal update variant: %w", err)
	}

	return &existing, nil
}

func (s *productVariantService) Delete(id uuid.UUID) error {
	variant, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	count, err := s.repo.CountByProductId(*variant.ProductId)
	if err != nil {
		return err
	}

	if count == 0 {
		if err := s.productRepo.ResetVariantStateToFalse(*variant.ProductId); err != nil {
			return err
		}
	}

	return nil
}

func (s *productVariantService) FindById(id uuid.UUID) (*entity.ProductVariant, error) {
	variant, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (s *productVariantService) FindByProductId(productId uuid.UUID) ([]entity.ProductVariant, error) {
	return s.repo.FindByProductId(productId)
}

func (s *productVariantService) SetActive(id uuid.UUID, isActive bool) error {
	return s.repo.SetActive(id, isActive)
}

func (s *productVariantService) SetAvailable(id uuid.UUID, isAvailable bool) error {
	return s.repo.SetAvailable(id, isAvailable)
}
