package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ProductService interface {
	Create(req request.ProductCreate) error
	Update(id int, req request.ProductUpdate) error
	Delete(id int) error
	FindById(id int) (response.ProductResponse, error)
	FindAll() ([]response.ProductResponse, error)
}

type productService struct {
	ProductRepo        repository.ProductRepository
	ProductVariantRepo repository.ProductVariantRepository
	Validate           *validator.Validate
}

func NewProductService(productRepo repository.ProductRepository, variantRepo repository.ProductVariantRepository, validate *validator.Validate) ProductService {
	return &productService{
		ProductRepo:        productRepo,
		ProductVariantRepo: variantRepo,
		Validate:           validate,
	}
}

func (s *productService) Create(req request.ProductCreate) error {
	err := s.Validate.Struct(req)
	if err != nil {
		return err
	}

	product := entity.Product{
		BusinessId:        req.BusinessId,
		ProductCategoryId: req.ProductCategoryId,
		Name:              req.Name,
		Description:       req.Description,
		Image:             req.Image,
		BasePrice:         req.BasePrice,
		Discount:          req.Discount,
		Promo:             req.Promo,
		Stock:             req.Stock,
		FinalPrice:        req.BasePrice - req.Discount - req.Promo,
		SKU:               req.SKU,
		HasVariant:        len(req.Variants) > 0,
	}
	product.Prepare()

	err = s.ProductRepo.InsertProduct(&product)
	if err != nil {
		return err
	}
	// Handle variants
	for _, v := range req.Variants {
		variant := entity.ProductVariant{
			BusinessId:  v.BusinessId,
			ProductId:   product.Id,
			Name:        v.Name,
			Image:       v.Image,
			BasePrice:   v.BasePrice,
			Discount:    v.Discount,
			Promo:       v.Promo,
			FinalPrice:  v.BasePrice - v.Discount - v.Promo,
			SKU:         v.SKU,
			Stock:       v.Stock,
			IsAvailable: true,
			IsActive:    true,
		}
		variant.Prepare()
		_ = s.ProductVariantRepo.Create(&variant)
	}

	return nil
}

func (s *productService) Update(id int, req request.ProductUpdate) error {
	err := s.Validate.Struct(req)
	if err != nil {
		return err
	}

	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Image = req.Image
	product.BasePrice = req.BasePrice
	product.Discount = req.Discount
	product.Promo = req.Promo
	product.Stock = req.Stock
	product.FinalPrice = req.BasePrice - req.Discount - req.Promo
	product.SKU = req.SKU
	product.IsAvailable = req.IsAvailable
	product.IsActive = req.IsActive
	product.HasVariant = len(req.Variants) > 0

	err = s.ProductRepo.UpdateProduct(&product)
	if err != nil {
		return err
	}

	// Hapus semua variant lama
	_ = s.ProductVariantRepo.DeleteByProductId(product.Id)

	// Tambahkan ulang variant baru
	for _, v := range req.Variants {
		variant := entity.ProductVariant{
			BusinessId:  v.BusinessId,
			ProductId:   product.Id,
			Name:        v.Name,
			Image:       v.Image,
			BasePrice:   v.BasePrice,
			Discount:    v.Discount,
			Promo:       v.Promo,
			FinalPrice:  v.BasePrice - v.Discount - v.Promo,
			SKU:         v.SKU,
			Stock:       v.Stock,
			IsAvailable: v.IsAvailable,
			IsActive:    v.IsActive,
		}
		variant.Prepare()
		_ = s.ProductVariantRepo.Create(&variant)
	}

	return nil
}

func (s *productService) Delete(id int) error {
	return s.ProductRepo.Delete(id)
}

func (s *productService) FindById(id int) (response.ProductResponse, error) {
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, err
	}
	return mapProductToResponse(product), nil
}

func (s *productService) FindAll() ([]response.ProductResponse, error) {
	products, err := s.ProductRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.ProductResponse
	for _, p := range products {
		responses = append(responses, mapProductToResponse(p))
	}
	return responses, nil
}

// Helper function to map entity to response
func mapProductToResponse(product entity.Product) response.ProductResponse {
	var variants []response.ProductVariantResponse
	for _, variant := range product.Variants {
		variants = append(variants, response.ProductVariantResponse{
			Id:        variant.Id,
			Name:      variant.Name,
			BasePrice: variant.BasePrice,
			SKU:       variant.SKU,
		})
	}

	var categoryRes *response.ProductCategoryResponse
	if product.ProductCategory.Id != 0 {
		categoryRes = &response.ProductCategoryResponse{
			Id:   product.ProductCategory.Id,
			Name: product.ProductCategory.Name,
		}
	}

	return response.ProductResponse{
		Id:                product.Id,
		Name:              product.Name,
		Description:       product.Description,
		Image:             product.Image,
		BasePrice:         product.BasePrice,
		FinalPrice:        product.FinalPrice,
		Discount:          product.Discount,
		Promo:             product.Promo,
		SKU:               product.SKU,
		Stock:             product.Stock,
		IsAvailable:       product.IsAvailable,
		IsActive:          product.IsActive,
		HasVariant:        product.HasVariant,
		Variants:          variants,
		ProductCategoryId: product.ProductCategoryId,
		ProductCategory:   categoryRes,
	}
}
