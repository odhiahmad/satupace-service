package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/redis/go-redis/v9"
)

type AutocompleteProduct struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type ProductService interface {
	Create(req request.ProductRequest) (response.ProductResponse, error)
	Update(id int, req request.ProductRequest) (response.ProductResponse, error)
	Delete(id int) error
	FindById(id int) (response.ProductResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error)
	SearchProducts(businessId int, search string, limit int) ([]response.ProductResponse, int64, error)
	SetActive(id int, isActive bool) error
	SetAvailable(id int, isAvailable bool) error
	UpdateImage(id int, base64Image string) (response.ProductResponse, error)
}

type productService struct {
	ProductRepo        repository.ProductRepository
	ProductVariantRepo repository.ProductVariantRepository
	Validate           *validator.Validate
	Redis              *redis.Client
}

func NewProductService(productRepo repository.ProductRepository, variantRepo repository.ProductVariantRepository, validate *validator.Validate, redis *redis.Client) ProductService {
	return &productService{
		ProductRepo:        productRepo,
		ProductVariantRepo: variantRepo,
		Validate:           validate,
		Redis:              redis,
	}
}

func (s *productService) Create(req request.ProductRequest) (response.ProductResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return response.ProductResponse{}, err
	}

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
		trackStock = true
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
		createdProduct, err := txRepo.Create(product)
		if err != nil {
			return fmt.Errorf("gagal menyimpan produk: %w", err)
		}
		product = createdProduct

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

		return nil
	})

	if err != nil {
		return response.ProductResponse{}, err
	}

	createdProduct, err := s.ProductRepo.FindById(product.Id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal mengambil data produk setelah simpan: %w", err)
	}

	if createdProduct.Image != nil {
		go func() {
			err := helper.AddProductToAutocomplete(s.Redis, req.BusinessId, product.Id, product.Name, *createdProduct.Image)
			if err != nil {
				fmt.Printf("gagal menambahkan autocomplete Redis: %v\n", err)
			}
		}()
	}

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
		trackStock = true
	}

	oldName := product.Name
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

	_ = helper.UpdateProductAutocomplete(s.Redis, product.BusinessId, oldName, product.Name, product.Id, *product.Image)

	productResponse := helper.MapProductToResponse(updatedProduct)

	return productResponse, nil
}

func (s *productService) Delete(id int) error {
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return err
	}

	_ = s.ProductVariantRepo.DeleteByProductId(id)
	err = s.ProductRepo.Delete(id)
	if err != nil {
		return err
	}

	if err := helper.DeleteProductFromAutocomplete(s.Redis, product.BusinessId, product.Name); err != nil {
		fmt.Printf("gagal menghapus autocomplete produk: %v\n", err)
	}

	for _, variant := range product.Variants {
		if err := helper.DeleteProductFromAutocomplete(s.Redis, product.BusinessId, variant.Name); err != nil {
			fmt.Printf("gagal menghapus autocomplete variant: %v\n", err)
		}
	}

	return nil
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
	products, total, err := s.ProductRepo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.ProductResponse
	for _, p := range products {
		result = append(result, helper.MapProductToResponse(p))
	}
	return result, total, nil
}

func (s *productService) SearchProducts(businessId int, search string, limit int) ([]response.ProductResponse, int64, error) {
	// Ambil hasil autocomplete dari Redis
	results, err := helper.GetProductAutocomplete(s.Redis, businessId, search, int64(limit))
	if err != nil {
		return nil, 0, err
	}
	if len(results) == 0 {
		return nil, 0, nil
	}

	// Parse hasil JSON dan ambil ID produk
	var productIDs []int
	autocompleteMap := make(map[int]response.ProductResponse)

	for _, product := range results {
		productIDs = append(productIDs, product.Id)
		autocompleteMap[product.Id] = product
	}

	// Ambil produk dari database
	products, err := s.ProductRepo.FindByIds(businessId, productIDs)
	if err != nil {
		return nil, 0, err
	}

	// Build hasil akhir, tambahkan image dari Redis jika ada
	var finalResults []response.ProductResponse
	for _, p := range products {
		res := helper.MapProductToResponse(p)

		if fromRedis, ok := autocompleteMap[p.Id]; ok && fromRedis.Image != nil && *fromRedis.Image != "" {
			res.Image = fromRedis.Image // override jika ada image dari Redis
		}

		finalResults = append(finalResults, res)
	}

	return finalResults, int64(len(finalResults)), nil
}

func (s *productService) SearchProductsRedisOnly(businessId int, search string, limit int) ([]response.ProductResponse, int64, error) {
	// Ambil hasil autocomplete dari Redis
	results, err := helper.GetProductAutocomplete(s.Redis, businessId, search, int64(limit))
	if err != nil {
		return nil, 0, err
	}

	// Jika tidak ada hasil, langsung kembalikan kosong
	if len(results) == 0 {
		return nil, 0, nil
	}

	// Tidak perlu ambil ke database, cukup return hasil dari Redis
	return results, int64(len(results)), nil
}

func (s *productService) UpdateImage(id int, base64Image string) (response.ProductResponse, error) {
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("produk tidak ditemukan: %w", err)
	}

	var oldImageURL *string = product.Image

	newImageURL, err := helper.UploadBase64ToCloudinary(base64Image, "product")
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal upload gambar baru: %w", err)
	}
	product.Image = &newImageURL

	updatedProduct, err := s.ProductRepo.UpdateAll(&product)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal update produk: %w", err)
	}

	if oldImageURL != nil {
		publicID, err := helper.ExtractPublicIDFromURL(*oldImageURL)
		if err == nil {
			_ = helper.DeleteFromCloudinary(publicID)
		}
	}

	return helper.MapProductToResponse(updatedProduct), nil
}
