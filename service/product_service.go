package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ProductService interface {
	Create(req request.ProductCreate) error
	Update(id int, req request.ProductUpdate) error
	Delete(id int) error
	FindById(id int) (response.ProductResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error) // <- Tambahan
}

type productService struct {
	ProductRepo        repository.ProductRepository
	ProductVariantRepo repository.ProductVariantRepository
	ProductPromoRepo   repository.ProductPromoRepository
	Validate           *validator.Validate
}

func NewProductService(productRepo repository.ProductRepository, variantRepo repository.ProductVariantRepository, promoRepo repository.ProductPromoRepository, validate *validator.Validate) ProductService {
	return &productService{
		ProductRepo:        productRepo,
		ProductVariantRepo: variantRepo,
		ProductPromoRepo:   promoRepo,
		Validate:           validate,
	}
}

func (s *productService) Create(req request.ProductCreate) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	categoryId := helper.IntOrDefault(req.ProductCategoryId, 0)
	basePrice := helper.Float64OrDefault(req.BasePrice, 0.0)
	sku := helper.StringOrDefault(req.SKU, "")
	stock := helper.IntOrDefault(req.Stock, 0)

	product := entity.Product{
		BusinessId:        req.BusinessId,
		ProductCategoryId: categoryId,
		Name:              req.Name,
		Description:       req.Description,
		Image:             req.Image,
		BasePrice:         basePrice,
		SKU:               sku,
		Stock:             stock,
		HasVariant:        len(req.Variants) > 0,
		IsAvailable:       true,
		IsActive:          true,
	}

	// ✅ simpan dengan pointer agar ID terisi
	err := s.ProductRepo.Create(&product)
	if err != nil {
		return err
	}

	// ✅ Gunakan product.Id yang sudah terisi
	for _, v := range req.Variants {
		variant := entity.ProductVariant{
			BusinessId:  req.BusinessId,
			ProductId:   product.Id,
			Name:        v.Name,
			Image:       v.Image,
			BasePrice:   v.BasePrice,
			SKU:         v.SKU,
			Stock:       v.Stock,
			IsAvailable: true,
			IsActive:    true,
		}
		_ = s.ProductVariantRepo.Create(&variant)
	}

	var productPromos []entity.ProductPromo
	for _, promoId := range req.PromoIds {
		productPromos = append(productPromos, entity.ProductPromo{
			BusinessId: req.BusinessId,
			ProductId:  product.Id,
			PromoId:    promoId,
		})
	}

	if len(productPromos) > 0 {
		err := s.ProductPromoRepo.CreateMany(productPromos)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *productService) Update(id int, req request.ProductUpdate) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return err
	}

	categoryId := helper.IntOrDefault(req.ProductCategoryId, 0)
	basePrice := helper.Float64OrDefault(req.BasePrice, 0.0)
	sku := helper.StringOrDefault(req.SKU, "")
	stock := helper.IntOrDefault(req.Stock, 0)

	product.ProductCategoryId = categoryId
	product.Name = req.Name
	product.Description = req.Description
	product.Image = req.Image
	product.BasePrice = basePrice
	product.SKU = sku
	product.Stock = stock
	product.IsAvailable = req.IsAvailable
	product.IsActive = req.IsActive
	product.HasVariant = len(req.Variants) > 0
	product.TaxId = req.TaxId
	product.DiscountId = req.DiscountId
	product.UnitId = req.UnitId

	updatedProduct, err := s.ProductRepo.Update(product)
	if err != nil {
		return err
	}
	// Kalau ingin tetap pakai updatedProduct
	product = updatedProduct

	// Hapus & update Variant
	_ = s.ProductVariantRepo.DeleteByProductId(product.Id)
	for _, v := range req.Variants {
		variant := entity.ProductVariant{
			ProductId:   product.Id,
			Name:        v.Name,
			Image:       v.Image,
			BasePrice:   v.BasePrice,
			SKU:         v.SKU,
			Stock:       v.Stock,
			IsAvailable: v.IsAvailable,
			IsActive:    v.IsActive,
		}
		_ = s.ProductVariantRepo.Create(&variant)
	}

	// Hapus & tambahkan ulang ProductPromo
	_ = s.ProductPromoRepo.DeleteByProductId(product.Id)

	var promos []entity.ProductPromo
	for _, promoId := range req.PromoIds {
		promos = append(promos, entity.ProductPromo{
			ProductId: product.Id,
			PromoId:   promoId,
		})
	}

	_ = s.ProductPromoRepo.CreateMany(promos)

	return nil
}

func (s *productService) Delete(id int) error {
	_ = s.ProductVariantRepo.DeleteByProductId(id)
	_ = s.ProductPromoRepo.DeleteByProductId(id)
	return s.ProductRepo.Delete(id)
}

func (s *productService) FindById(id int) (response.ProductResponse, error) {
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, err
	}
	return mapProductToResponse(product), nil
}

func (s *productService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error) {
	products, total, err := s.ProductRepo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.ProductResponse
	for _, product := range products {
		result = append(result, mapProductToResponse(product))
	}

	return result, total, nil
}

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

	var promos []response.ProductPromoResponse
	for _, promo := range product.ProductPromos {
		if promo.Promo.Id != 0 {
			promos = append(promos, response.ProductPromoResponse{
				PromoId:   promo.PromoId,
				ProductId: promo.ProductId,
			})
		}
	}

	var categoryRes *response.ProductCategoryResponse
	if product.ProductCategory != nil && product.ProductCategory.Id != 0 {
		categoryRes = &response.ProductCategoryResponse{
			Id:   product.ProductCategory.Id,
			Name: product.ProductCategory.Name,
		}
	}

	var taxRes *response.TaxResponse
	if product.Tax != nil {
		taxRes = &response.TaxResponse{
			Id:     product.Tax.Id,
			Name:   product.Tax.Name,
			Amount: product.Tax.Amount,
		}
	}

	var discountRes *response.DiscountResponse
	if product.Discount != nil {
		discountRes = &response.DiscountResponse{
			Id:     product.Discount.Id,
			Name:   product.Discount.Name,
			Amount: product.Discount.Amount,
		}
	}

	var unitRes *response.ProductUnitResponse
	if product.Unit != nil {
		unitRes = &response.ProductUnitResponse{
			Id:         product.Unit.Id,
			Name:       product.Unit.Name,
			Alias:      product.Unit.Alias,
			Multiplier: product.Unit.Multiplier,
		}
	}

	return response.ProductResponse{
		Id:                product.Id,
		Name:              product.Name,
		Description:       product.Description,
		Image:             product.Image,
		BasePrice:         product.BasePrice,
		FinalPrice:        product.FinalPrice,
		SKU:               product.SKU,
		Stock:             product.Stock,
		IsAvailable:       product.IsAvailable,
		IsActive:          product.IsActive,
		HasVariant:        product.HasVariant,
		Variants:          variants,
		ProductCategoryId: product.ProductCategoryId,
		ProductCategory:   categoryRes,
		Tax:               taxRes,
		Discount:          discountRes,
		Unit:              unitRes,
		Promos:            promos,
	}
}
