package service

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/redis/go-redis/v9"
)

type ProductService interface {
	Create(req request.ProductCreate) error
	Update(id int, req request.ProductUpdate) (*entity.Product, error)
	Delete(id int) error
	FindById(id int) (response.ProductResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error)
	SetActive(id int, isActive bool) error
	SetAvailable(id int, isAvailable bool) error
}

type productService struct {
	ProductRepo        repository.ProductRepository
	ProductPromoRepo   repository.ProductPromoRepository
	ProductVariantRepo repository.ProductVariantRepository
	PromoRepo          repository.PromoRepository
	Validate           *validator.Validate
	Redis              *redis.Client
}

func NewProductService(productRepo repository.ProductRepository, productPromoRepo repository.ProductPromoRepository, promoRepo repository.PromoRepository, variantRepo repository.ProductVariantRepository, validate *validator.Validate, redis *redis.Client) ProductService {
	return &productService{
		ProductRepo:        productRepo,
		ProductVariantRepo: variantRepo,
		ProductPromoRepo:   productPromoRepo,
		PromoRepo:          promoRepo,
		Validate:           validate,
		Redis:              redis,
	}
}

func (s *productService) Create(req request.ProductCreate) error {
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

	categoryId := helper.IntOrDefault(req.ProductCategoryId, 0)
	basePrice := helper.Float64OrDefault(req.BasePrice, 0.0)
	stock := helper.IntOrDefault(req.Stock, 0)
	sku := helper.StringOrDefault(req.SKU, "")
	trackStock := stock == 0

	if sku == "" {
		sku = helper.GenerateSKU(req.Name)
	}

	product := entity.Product{
		BusinessId:        req.BusinessId,
		ProductCategoryId: categoryId,
		Name:              req.Name,
		Description:       req.Description,
		Image:             imageURL,
		MinimumSales:      req.MinimumSales,
		BasePrice:         helper.Float64Ptr(basePrice),
		SKU:               helper.StringPtr(sku),
		Stock:             helper.IntPtr(stock),
		HasVariant:        len(req.Variants) > 0,
		TaxId:             req.TaxId,
		DiscountId:        req.DiscountId,
		UnitId:            req.UnitId,
		TrackStock:        trackStock,
		IsAvailable:       true,
		IsActive:          true,
	}

	// Jika punya variant, kosongkan harga dan stock di induk
	if product.HasVariant {
		product.BasePrice = nil
		product.Stock = nil
		product.SKU = nil
	}

	// Mulai transaksi
	return s.ProductRepo.WithTransaction(func(txRepo repository.ProductRepository) error {
		// Simpan product utama
		if err := txRepo.Create(&product); err != nil {
			return fmt.Errorf("gagal menyimpan produk: %w", err)
		}

		if err := helper.IndexProductToElastic(&product); err != nil {
			log.Printf("gagal mengindeks produk ke Elasticsearch: %v", err)
		}

		// Simpan variants
		for _, v := range req.Variants {
			skuVariant := v.SKU
			trackStockVariant := v.Stock == 0

			if skuVariant == "" {
				skuVariant = helper.GenerateSKU(v.Name)
			}

			variant := entity.ProductVariant{
				BusinessId:  req.BusinessId,
				ProductId:   product.Id,
				Name:        v.Name,
				BasePrice:   v.BasePrice,
				SKU:         skuVariant,
				Stock:       v.Stock,
				TrackStock:  trackStockVariant,
				IsAvailable: true,
				IsActive:    true,
			}

			if err := s.ProductVariantRepo.CreateWithTx(txRepo, &variant); err != nil {
				return fmt.Errorf("gagal menyimpan variant '%s': %w", v.Name, err)
			}
		}

		// Simpan relasi promo untuk produk utama
		if len(req.PromoIds) > 0 {
			var promos []entity.ProductPromo
			for _, promoId := range req.PromoIds {
				promos = append(promos, entity.ProductPromo{
					BusinessId: req.BusinessId,
					ProductId:  &product.Id,
					PromoId:    promoId,
				})
			}
			if err := s.ProductPromoRepo.CreateManyWithTx(txRepo, promos); err != nil {
				return fmt.Errorf("gagal menyimpan relasi promo untuk produk: %w", err)
			}
		}

		return nil
	})
}

func (s *productService) Update(id int, req request.ProductUpdate) (*entity.Product, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	hasVariants := len(req.Variants) > 0
	trackStock := *req.Stock == 0

	product.ProductCategoryId = *req.ProductCategoryId
	product.Name = req.Name
	product.Description = req.Description
	product.HasVariant = hasVariants
	product.IsAvailable = true
	product.IsActive = true
	product.TaxId = req.TaxId
	product.UnitId = req.UnitId
	product.TrackStock = trackStock
	product.DiscountId = req.DiscountId
	product.MinimumSales = req.MinimumSales

	// Jika punya variant, kosongkan harga dan stock di induk
	if hasVariants {
		product.BasePrice = nil
		product.Stock = nil
		product.SKU = nil
	} else {
		product.BasePrice = req.BasePrice
		product.Stock = req.Stock
		product.SKU = req.SKU
	}

	updatedProduct, err := s.ProductRepo.Update(product)
	if err != nil {
		return nil, err
	}

	if err := helper.IndexProductToElastic(&product); err != nil {
		log.Printf("gagal mengindeks produk ke Elasticsearch: %v", err)
	}

	_ = s.ProductPromoRepo.DeleteByProductId(product.Id)

	if len(req.PromoIds) > 0 {
		var promos []entity.ProductPromo
		for _, promoId := range req.PromoIds {
			exists, err := s.PromoRepo.Exists(promoId)
			if err != nil {
				return nil, fmt.Errorf("gagal cek promo: %w", err)
			}
			if !exists {
				return nil, fmt.Errorf("promo dengan ID %d tidak ditemukan", promoId)
			}

			promos = append(promos, entity.ProductPromo{
				BusinessId: product.BusinessId,
				ProductId:  &product.Id,
				PromoId:    promoId,
			})
		}
		if err := s.ProductPromoRepo.CreateMany(promos); err != nil {
			return nil, fmt.Errorf("gagal menyimpan relasi promo baru: %w", err)
		}
	}

	return &updatedProduct, nil
}

func (s *productService) Delete(id int) error {
	_ = s.ProductVariantRepo.DeleteByProductId(id)
	return s.ProductRepo.Delete(id)
}

func (s *productService) SetActive(id int, isActive bool) error {
	return s.ProductRepo.SetActive(id, isActive)
}

func (s *productService) SetAvailable(id int, isAvailable bool) error {
	return s.ProductRepo.SetAvailable(id, isAvailable)
}

func (s *productService) FindById(id int) (response.ProductResponse, error) {
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, err
	}
	res := helper.MapProductToResponse(product)
	return res, nil
}

func (s *productService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error) {
	var products []entity.Product
	var total int64
	var err error

	products, total, err = s.ProductRepo.FindWithPagination(businessId, pagination)

	if err != nil {
		return nil, 0, err
	}

	var result []response.ProductResponse
	for _, product := range products {
		result = append(result, helper.MapProductToResponse(product))
	}
	return result, total, nil
}
