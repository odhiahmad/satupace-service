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
	Create(req request.ProductRequest) (response.ProductResponse, error)
	Update(id int, req request.ProductRequest) (response.ProductResponse, error)
	Delete(id int) error
	FindById(id int) (response.ProductResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error)
	SetActive(id int, isActive bool) error
	SetAvailable(id int, isAvailable bool) error
	UpdateImage(id int, base64Image string) (response.ProductResponse, error)
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

func (s *productService) Create(req request.ProductRequest) (response.ProductResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return response.ProductResponse{}, err
	}

	// Upload gambar produk utama
	var imageURL *string
	if req.Image != nil && *req.Image != "" {
		url, err := helper.UploadBase64ToCloudinary(*req.Image, "product")
		if err != nil {
			return response.ProductResponse{}, fmt.Errorf("gagal upload gambar produk: %w", err)
		}
		imageURL = &url
	}

	sku := helper.StringOrDefault(req.SKU, "")

	var trackStock bool
	if req.Stock != nil {
		trackStock = *req.Stock == 0
	} else {
		trackStock = true // atau false tergantung default logika kamu
	}

	if sku == "" {
		sku = helper.GenerateSKU(req.Name)
	}

	product := entity.Product{
		BusinessId:   req.BusinessId,
		CategoryId:   *req.CategoryId,
		Name:         req.Name,
		Description:  req.Description,
		Image:        imageURL,
		MinimumSales: req.MinimumSales,
		BasePrice:    req.BasePrice,
		SellPrice:    req.SellPrice,
		SKU:          helper.StringPtr(sku),
		Stock:        req.Stock,
		HasVariant:   len(req.Variants) > 0,
		BrandId:      req.BrandId,
		TaxId:        req.TaxId,
		DiscountId:   req.DiscountId,
		UnitId:       req.UnitId,
		TrackStock:   trackStock,
		IsAvailable:  true,
		IsActive:     true,
	}

	if product.HasVariant {
		product.BasePrice = nil
		product.Stock = nil
		product.SKU = nil
	}

	err := s.ProductRepo.WithTransaction(func(txRepo repository.ProductRepository) error {
		// Simpan produk
		createdProduct, err := txRepo.Create(product)
		if err != nil {
			return fmt.Errorf("gagal menyimpan produk: %w", err)
		}
		product = createdProduct // assign balik ke variabel utama jika diperlukan

		if err := helper.IndexProductToElastic(&product); err != nil {
			log.Printf("gagal mengindeks produk ke Elasticsearch: %v", err)
		}

		// Simpan variants
		for _, v := range req.Variants {
			skuVariant := v.SKU
			if skuVariant == "" {
				skuVariant = helper.GenerateSKU(v.Name)
			}

			variant := entity.ProductVariant{
				BusinessId:  req.BusinessId,
				ProductId:   product.Id,
				Name:        v.Name,
				BasePrice:   v.BasePrice,
				SellPrice:   v.SellPrice,
				SKU:         skuVariant,
				Stock:       v.Stock,
				TrackStock:  v.Stock == 0,
				IsAvailable: true,
				IsActive:    true,
			}

			if err := s.ProductVariantRepo.CreateWithTx(txRepo, &variant); err != nil {
				return fmt.Errorf("gagal menyimpan variant '%s': %w", v.Name, err)
			}
		}

		// Simpan relasi promo
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
				return fmt.Errorf("gagal menyimpan relasi promo: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return response.ProductResponse{}, err
	}

	// Ambil ulang produk yang sudah di-preload relasinya
	createdProduct, err := s.ProductRepo.FindById(product.Id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal mengambil data produk setelah simpan: %w", err)
	}

	// Mapping ke response
	productResponse := helper.MapProductToResponse(createdProduct)
	return productResponse, nil
}

func (s *productService) Update(id int, req request.ProductRequest) (response.ProductResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return response.ProductResponse{}, err
	}

	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, err
	}

	hasVariants := len(req.Variants) > 0

	var trackStock bool
	if req.Stock != nil {
		trackStock = *req.Stock == 0
	} else {
		trackStock = true // atau false tergantung default logika kamu
	}

	fmt.Println("TrackStock:", trackStock)

	product.CategoryId = *req.CategoryId
	product.Name = req.Name
	product.Description = req.Description
	product.HasVariant = hasVariants
	product.IsAvailable = true
	product.IsActive = true
	product.BrandId = req.BrandId
	product.TaxId = req.TaxId
	product.UnitId = req.UnitId
	product.TrackStock = trackStock
	product.DiscountId = req.DiscountId
	product.MinimumSales = req.MinimumSales

	// Jika punya variant, kosongkan harga dan stock di induk
	if hasVariants {
		product.BasePrice = nil
		product.SellPrice = nil
		product.Stock = nil
		product.SKU = nil
	} else {
		product.BasePrice = req.BasePrice
		product.SellPrice = req.SellPrice
		product.Stock = req.Stock
		product.SKU = req.SKU
	}

	updatedProduct, err := s.ProductRepo.Update(product)
	if err != nil {
		return response.ProductResponse{}, err
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
				return response.ProductResponse{}, fmt.Errorf("gagal cek promo: %w", err)
			}
			if !exists {
				return response.ProductResponse{}, fmt.Errorf("promo dengan ID %d tidak ditemukan", promoId)
			}

			promos = append(promos, entity.ProductPromo{
				BusinessId: product.BusinessId,
				ProductId:  &product.Id,
				PromoId:    promoId,
			})
		}
		if err := s.ProductPromoRepo.CreateMany(promos); err != nil {
			return response.ProductResponse{}, fmt.Errorf("gagal menyimpan relasi promo baru: %w", err)
		}
	}

	productResponse := helper.MapProductToResponse(updatedProduct)

	return productResponse, nil
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

func (s *productService) UpdateImage(id int, base64Image string) (response.ProductResponse, error) {
	// Cari produk
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("produk tidak ditemukan: %w", err)
	}

	// Simpan URL gambar lama
	var oldImageURL *string = product.Image

	// Upload gambar baru ke Cloudinary
	newImageURL, err := helper.UploadBase64ToCloudinary(base64Image, "product")
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal upload gambar baru: %w", err)
	}
	product.Image = &newImageURL

	// Update produk ke DB
	updatedProduct, err := s.ProductRepo.UpdateAll(&product)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal update produk: %w", err)
	}

	// Hapus gambar lama dari Cloudinary jika ada
	if oldImageURL != nil {
		publicID, err := helper.ExtractPublicIDFromURL(*oldImageURL)
		if err == nil {
			_ = helper.DeleteFromCloudinary(publicID)
		}
	}

	return helper.MapProductToResponse(updatedProduct), nil
}
