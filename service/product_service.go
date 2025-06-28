package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	ProductVariantRepo repository.ProductVariantRepository
	ProductPromoRepo   repository.ProductPromoRepository
	Validate           *validator.Validate
	Redis              *redis.Client
}

func NewProductService(productRepo repository.ProductRepository, variantRepo repository.ProductVariantRepository, promoRepo repository.ProductPromoRepository, validate *validator.Validate, redis *redis.Client) ProductService {
	return &productService{
		ProductRepo:        productRepo,
		ProductVariantRepo: variantRepo,
		ProductPromoRepo:   promoRepo,
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
	if sku == "" {
		sku = helper.GenerateSKU(req.Name)
	}

	product := entity.Product{
		BusinessId:        req.BusinessId,
		ProductCategoryId: categoryId,
		Name:              req.Name,
		Description:       req.Description,
		Image:             imageURL,
		BasePrice:         helper.Float64Ptr(basePrice),
		SKU:               helper.StringPtr(sku),
		Stock:             helper.IntPtr(stock),
		HasVariant:        len(req.Variants) > 0,
		TaxId:             req.TaxId,
		DiscountId:        req.DiscountId,
		UnitId:            req.UnitId,
		IsAvailable:       true,
		IsActive:          true,
	}

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

		// Simpan variant dan promosinya
		for _, v := range req.Variants {
			var variantImageURL *string
			if v.Image != nil && *v.Image != "" {
				url, err := helper.UploadBase64ToCloudinary(*v.Image, "variant")
				if err != nil {
					return fmt.Errorf("gagal upload gambar variant '%s': %w", v.Name, err)
				}
				variantImageURL = &url
			}

			skuVariant := v.SKU
			if skuVariant == "" {
				skuVariant = helper.GenerateSKU(v.Name)
			}

			variant := entity.ProductVariant{
				BusinessId:  req.BusinessId,
				ProductId:   product.Id,
				Name:        v.Name,
				Image:       variantImageURL,
				BasePrice:   v.BasePrice,
				SKU:         skuVariant,
				Stock:       v.Stock,
				TaxId:       v.TaxId,
				DiscountId:  v.DiscountId,
				UnitId:      v.UnitId,
				IsAvailable: true,
				IsActive:    true,
			}

			if err := s.ProductVariantRepo.CreateWithTx(txRepo, &variant); err != nil {
				return fmt.Errorf("gagal menyimpan variant '%s': %w", v.Name, err)
			}

			// Simpan relasi promo untuk variant
			if len(v.PromoIds) > 0 {
				var promos []entity.ProductPromo
				for _, promoId := range v.PromoIds {
					promos = append(promos, entity.ProductPromo{
						BusinessId:       req.BusinessId,
						ProductVariantId: &variant.Id,
						PromoId:          promoId,
					})
				}
				if err := s.ProductPromoRepo.CreateManyWithTx(txRepo, promos); err != nil {
					return fmt.Errorf("gagal menyimpan relasi promo untuk variant '%s': %w", v.Name, err)
				}
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

	var imageURL *string
	if req.Image != nil && *req.Image != "" {
		url, err := helper.UploadBase64ToCloudinary(*req.Image, "product")
		if err != nil {
			return nil, fmt.Errorf("gagal upload gambar produk: %w", err)
		}
		imageURL = &url
	}

	categoryId := helper.IntOrDefault(req.ProductCategoryId, 0)
	basePrice := helper.Float64OrDefault(req.BasePrice, 0.0)
	stock := helper.IntOrDefault(req.Stock, 0)
	sku := helper.StringOrDefault(req.SKU, "")
	if sku == "" {
		sku = helper.GenerateSKU(req.Name)
	}

	hasVariants := len(req.Variants) > 0

	product.ProductCategoryId = categoryId
	product.Name = req.Name
	product.Description = req.Description
	product.Image = imageURL
	product.HasVariant = hasVariants
	product.IsAvailable = true
	product.IsActive = true
	product.TaxId = req.TaxId
	product.DiscountId = req.DiscountId
	product.UnitId = req.UnitId

	// jika punya variant, kosongkan harga dan stock di induk
	if hasVariants {
		product.BasePrice = nil
		product.Stock = nil
		product.SKU = nil
	} else {
		product.BasePrice = helper.Float64Ptr(basePrice)
		product.Stock = helper.IntPtr(stock)
		product.SKU = helper.StringPtr(sku)
	}

	updatedProduct, err := s.ProductRepo.Update(product)
	if err != nil {
		return nil, err
	}

	_ = s.ProductPromoRepo.DeleteByProductId(product.Id)

	if len(req.PromoIds) > 0 {
		var promos []entity.ProductPromo
		for _, promoId := range req.PromoIds {
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
	_ = s.ProductPromoRepo.DeleteByProductId(id)
	return s.ProductRepo.Delete(id)
}

func (s *productService) SetActive(id int, isActive bool) error {
	return s.ProductRepo.SetActive(id, isActive)
}

func (s *productService) SetAvailable(id int, isAvailable bool) error {
	return s.ProductRepo.SetAvailable(id, isAvailable)
}

func (s *productService) FindById(id int) (response.ProductResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%d", id)

	// 🔍 Coba ambil dari cache Redis
	if s.Redis != nil {
		cachedData, err := s.Redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedProduct response.ProductResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedProduct); err == nil {
				log.Println("✅ Product found in Redis cache")
				return cachedProduct, nil
			}
		}
	}

	// 🗄 Ambil dari DB jika tidak ditemukan atau gagal decode cache
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, err
	}

	res := mapProductToResponse(product)

	// 💾 Simpan ke Redis
	if s.Redis != nil {
		if jsonData, err := json.Marshal(res); err == nil {
			err = s.Redis.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
			if err != nil {
				log.Println("❗️ Failed to cache product:", err)
			}
		}
	}

	return res, nil
}

func (s *productService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error) {
	ctx := context.Background()
	pong, err := s.Redis.Ping(ctx).Result()
	log.Println("🔗 Redis ping:", pong, err)
	cacheKey := fmt.Sprintf(
		"product:list:%d:%d:%d:%s:%s:%s",
		businessId,
		pagination.Page,
		pagination.Limit,
		pagination.SortBy,
		pagination.OrderBy,
		pagination.Search,
	)

	// Coba ambil dari cache Redis
	if s.Redis != nil {
		cachedData, err := s.Redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cached struct {
				Products []response.ProductResponse `json:"products"`
				Total    int64                      `json:"total"`
			}
			if err := json.Unmarshal([]byte(cachedData), &cached); err == nil {
				return cached.Products, cached.Total, nil
			}
		}
	}

	// Ambil dari DB jika tidak ada atau gagal cache
	products, total, err := s.ProductRepo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.ProductResponse
	for _, product := range products {
		result = append(result, mapProductToResponse(product))
	}

	// Simpan ke cache Redis
	if s.Redis != nil {
		cacheBody := struct {
			Products []response.ProductResponse `json:"products"`
			Total    int64                      `json:"total"`
		}{
			Products: result,
			Total:    total,
		}
		if jsonData, err := json.Marshal(cacheBody); err == nil {
			s.Redis.Set(ctx, cacheKey, jsonData, 5*time.Minute)
		}
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
			var requiredProducts []response.RequiredProductData
			for _, p := range promo.Promo.RequiredProducts {
				requiredProducts = append(requiredProducts, response.RequiredProductData{
					Id:   p.Id,
					Name: p.Name,
				})
			}

			promos = append(promos, response.ProductPromoResponse{
				Name:             promo.Promo.Name,
				Description:      helper.StringPtr(promo.Promo.Description),
				Amount:           promo.Promo.Amount,
				Type:             promo.Promo.Type,
				MinQuantity:      promo.Promo.MinQuantity,
				StartDate:        promo.Promo.StartDate,
				EndDate:          promo.Promo.EndDate,
				RequiredProducts: requiredProducts,
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
			Type:   product.Tax.Type,
		}
	}

	var discountRes *response.DiscountResponse
	if product.Discount != nil {
		discountRes = &response.DiscountResponse{
			Id:     product.Discount.Id,
			Name:   product.Discount.Name,
			Amount: product.Discount.Amount,
			Type:   product.Discount.Type,
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
		Id:              product.Id,
		Name:            product.Name,
		Description:     product.Description,
		Image:           product.Image,
		BasePrice:       helper.Float64OrDefault(product.BasePrice, 0.0),
		SKU:             helper.StringOrDefault(product.SKU, ""),
		Stock:           helper.IntOrDefault(product.Stock, 1),
		IsAvailable:     product.IsAvailable,
		IsActive:        product.IsActive,
		HasVariant:      product.HasVariant,
		Variants:        variants,
		ProductCategory: categoryRes,
		Tax:             taxRes,
		Discount:        discountRes,
		Unit:            unitRes,
		Promos:          promos,
	}
}
