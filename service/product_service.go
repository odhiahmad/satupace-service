package service

import (
	"fmt"
	"log"
	"strings"

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
	Update(id int, req request.ProductUpdateRequest) (response.ProductResponse, error)
	Delete(id int) error
	SearchProducts(businessId int, search string, limit int) ([]response.ProductResponse, int64, error)
	SearchProductsRedisOnly(businessId int, search string, limit int) ([]response.ProductSearchResponse, error)
	SetActive(id int, isActive bool) error
	SetAvailable(id int, isAvailable bool) error
	UpdateImage(id int, base64Image string) (response.ProductResponse, error)
	FindById(id int) (response.ProductResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.ProductResponse, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.ProductResponse, string, bool, error)
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
	falseVal := false

	if err := s.Validate.Struct(req); err != nil {
		return response.ProductResponse{}, err
	}

	hasVariant := len(req.Variants) > 0
	sku := req.SKU

	if hasVariant {
		skuMap := map[string]bool{}
		for i, v := range req.Variants {
			if v.SKU == nil || *v.SKU == "" {
				sGenerated := helper.GenerateSKU(strings.ToLower(req.Name))
				req.Variants[i].SKU = &sGenerated
			}

			if skuMap[*req.Variants[i].SKU] {
				return response.ProductResponse{}, fmt.Errorf("SKU varian duplikat di antara varian: %s", *req.Variants[i].SKU)
			}
			skuMap[*req.Variants[i].SKU] = true

			exist, err := s.ProductVariantRepo.IsSKUExist(*req.Variants[i].SKU, *req.BusinessId)
			if err != nil {
				return response.ProductResponse{}, fmt.Errorf("gagal cek SKU varian di database: %w", err)
			}
			if exist {
				return response.ProductResponse{}, fmt.Errorf("SKU varian sudah digunakan: %s", *req.Variants[i].SKU)
			}
		}
	} else {
		if sku == nil || *sku == "" {
			sGenerated := helper.GenerateSKU(strings.ToLower(req.Name))
			sku = &sGenerated
		}

		exist, err := s.ProductRepo.IsSKUExist(*sku, *req.BusinessId)
		if err != nil {
			return response.ProductResponse{}, fmt.Errorf("gagal cek SKU produk di database: %w", err)
		}
		if exist {
			return response.ProductResponse{}, fmt.Errorf("SKU produk sudah digunakan")
		}
	}

	product := entity.Product{
		BusinessId:   req.BusinessId,
		CategoryId:   req.CategoryId,
		Name:         strings.ToLower(req.Name),
		Description:  req.Description,
		Image:        nil,
		MinimumSales: req.MinimumSales,
		HasVariant:   hasVariant,
		BrandId:      req.BrandId,
		TaxId:        req.TaxId,
		DiscountId:   req.DiscountId,
		UnitId:       req.UnitId,
		IsAvailable:  req.IsAvailable,
		IsActive:     req.IsActive,
	}

	if hasVariant {
		product.BasePrice = nil
		product.SellPrice = nil
		product.Stock = nil
		product.SKU = nil
		product.TrackStock = &falseVal
		product.IgnoreStockCheck = &falseVal
	} else {
		product.BasePrice = req.BasePrice
		product.SellPrice = req.SellPrice
		product.Stock = req.Stock
		product.SKU = sku
		product.TrackStock = req.TrackStock
		product.IgnoreStockCheck = req.IgnoreStockCheck
	}

	err := s.ProductRepo.WithTransaction(func(txRepo repository.ProductRepository) error {
		createdProduct, err := txRepo.Create(product)
		if err != nil {
			return fmt.Errorf("gagal menyimpan produk: %w", err)
		}
		product = createdProduct

		if hasVariant {
			var variants []entity.ProductVariant
			for _, v := range req.Variants {
				variants = append(variants, entity.ProductVariant{
					BusinessId:       req.BusinessId,
					ProductId:        &product.Id,
					Name:             v.Name,
					Description:      v.Description,
					BasePrice:        v.BasePrice,
					SellPrice:        v.SellPrice,
					MinimumSales:     v.MinimumSales,
					SKU:              v.SKU,
					Stock:            v.Stock,
					TrackStock:       v.TrackStock,
					IgnoreStockCheck: req.IgnoreStockCheck,
					IsAvailable:      v.IsAvailable,
					IsActive:         v.IsActive,
				})
			}

			if err := s.ProductVariantRepo.CreateWithTx(txRepo, variants); err != nil {
				return fmt.Errorf("gagal menyimpan variant produk: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return response.ProductResponse{}, err
	}

	if req.Image != nil && *req.Image != "" {
		go func(productId int, businessId int, name string, imageBase64 string) {
			url, err := helper.UploadBase64ToCloudinary(imageBase64, "product")
			if err != nil {
				log.Printf("Gagal upload gambar produk (background): %v", err)
				return
			}

			err = s.ProductRepo.UpdateImage(productId, url)
			if err != nil {
				log.Printf("Gagal update gambar produk setelah upload: %v", err)
				return
			}

		}(product.Id, *req.BusinessId, strings.ToLower(req.Name), *req.Image)
	}

	createdProduct, err := s.ProductRepo.FindById(product.Id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal mengambil data produk setelah simpan: %w", err)
	}

	if err := helper.AddProductToAutocomplete(s.Redis, *req.BusinessId, product.Id, req.Name); err != nil {
		log.Printf("[Redis Autocomplete] Gagal menambahkan: %v", err)
	}

	return helper.MapProductToResponse(createdProduct), nil
}

func (s *productService) Update(id int, req request.ProductUpdateRequest) (response.ProductResponse, error) {
	falseVal := false

	if err := s.Validate.Struct(req); err != nil {
		return response.ProductResponse{}, err
	}

	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("produk tidak ditemukan: %w", err)
	}

	hasVariant := len(req.Variants) > 0
	sku := req.SKU
	oldName := product.Name

	if hasVariant {
		skuMap := map[string]bool{}
		for _, v := range req.Variants {
			if v.SKU == nil || *v.SKU == "" {
				continue
			}

			if skuMap[*v.SKU] {
				return response.ProductResponse{}, fmt.Errorf("SKU varian duplikat di antara varian: %s", *v.SKU)
			}
			skuMap[*v.SKU] = true

			exist, err := s.ProductVariantRepo.IsSKUExistExcept(*v.SKU, *req.BusinessId, v.Id)

			fmt.Println(v.SKU, v.Id, "sku dan id")
			if err != nil {
				return response.ProductResponse{}, fmt.Errorf("gagal cek SKU varian di database: %w", err)
			}
			if exist {
				return response.ProductResponse{}, fmt.Errorf("SKU varian sudah digunakan: %s", *v.SKU)
			}
		}
	} else {
		if sku == nil || *sku == "" {
			sGenerated := helper.GenerateSKU(strings.ToLower(req.Name))
			sku = &sGenerated
		}

		if product.SKU == nil || *product.SKU != *sku {
			exist, err := s.ProductRepo.IsSKUExistExcept(*sku, *req.BusinessId, id)
			if err != nil {
				return response.ProductResponse{}, fmt.Errorf("gagal cek SKU produk: %w", err)
			}
			if exist {
				return response.ProductResponse{}, fmt.Errorf("SKU produk sudah digunakan")
			}
		}
	}

	product.CategoryId = req.CategoryId
	product.Name = strings.ToLower(req.Name)
	product.Description = req.Description
	product.BrandId = req.BrandId
	product.TaxId = req.TaxId
	product.UnitId = req.UnitId
	product.DiscountId = req.DiscountId
	product.HasVariant = hasVariant
	product.IsAvailable = req.IsAvailable
	product.IsActive = req.IsActive

	if hasVariant {
		product.BasePrice = nil
		product.SellPrice = nil
		product.Stock = nil
		product.SKU = nil
		product.TrackStock = &falseVal
		product.IgnoreStockCheck = &falseVal
		product.MinimumSales = nil
	} else {
		product.BasePrice = req.BasePrice
		product.SellPrice = req.SellPrice
		product.Stock = req.Stock
		product.SKU = sku
		product.TrackStock = req.TrackStock
		product.IgnoreStockCheck = req.IgnoreStockCheck
		product.MinimumSales = req.MinimumSales
	}

	err = s.ProductRepo.WithTransaction(func(txRepo repository.ProductRepository) error {
		if _, err := txRepo.Update(product); err != nil {
			return fmt.Errorf("gagal update produk: %w", err)
		}

		if hasVariant {
			var updatedVariants []entity.ProductVariant
			for _, v := range req.Variants {
				existingVariant, err := s.ProductVariantRepo.FindById(v.Id)
				if err != nil {
					return fmt.Errorf("gagal ambil data variant: %w", err)
				}

				sku := existingVariant.SKU
				if v.SKU != nil {
					sku = v.SKU
				}

				updatedVariants = append(updatedVariants, entity.ProductVariant{
					Id:               v.Id,
					ProductId:        &product.Id,
					BusinessId:       req.BusinessId,
					Name:             v.Name,
					Description:      v.Description,
					BasePrice:        v.BasePrice,
					SellPrice:        v.SellPrice,
					MinimumSales:     v.MinimumSales,
					SKU:              sku,
					Stock:            v.Stock,
					TrackStock:       v.TrackStock,
					IgnoreStockCheck: v.IgnoreStockCheck,
					IsAvailable:      v.IsAvailable,
					IsActive:         v.IsActive,
				})
			}

			if err := s.ProductVariantRepo.UpdateWithTx(txRepo, updatedVariants); err != nil {
				return fmt.Errorf("gagal update variant produk: %w", err)
			}
		} else {
			if err := s.ProductVariantRepo.DeleteByProductId(product.Id); err != nil {
				return fmt.Errorf("gagal menghapus variant lama: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return response.ProductResponse{}, err
	}

	_ = helper.UpdateProductAutocomplete(s.Redis, *product.BusinessId, oldName, product.Name, product.Id)

	createdProduct, err := s.ProductRepo.FindById(product.Id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal mengambil data produk setelah update: %w", err)
	}

	return helper.MapProductToResponse(createdProduct), nil
}

func (s *productService) Delete(id int) error {
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return err
	}

	hasRelation, err := s.ProductRepo.HasRelation(id)
	if err != nil {
		return err
	}

	_ = s.ProductVariantRepo.DeleteByProductId(id)

	var deleteErr error
	if hasRelation {
		deleteErr = s.ProductRepo.SoftDelete(id)
	} else {
		deleteErr = s.ProductRepo.HardDelete(id)
	}
	if deleteErr != nil {
		return deleteErr
	}

	if err := helper.DeleteProductFromAutocomplete(s.Redis, *product.BusinessId, product.Name); err != nil {
		fmt.Printf("gagal menghapus autocomplete produk: %v\n", err)
	}
	for _, variant := range product.Variants {
		if err := helper.DeleteProductFromAutocomplete(s.Redis, *product.BusinessId, variant.Name); err != nil {
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

func (s *productService) SearchProducts(businessId int, search string, limit int) ([]response.ProductResponse, int64, error) {
	results, err := helper.GetProductAutocomplete(s.Redis, businessId, search, int64(limit))
	if err != nil {
		return nil, 0, err
	}
	if len(results) == 0 {
		return nil, 0, nil
	}

	var productIDs []int
	autocompleteMap := make(map[int]response.ProductResponse)

	for _, product := range results {
		productIDs = append(productIDs, product.Id)
		autocompleteMap[product.Id] = product
	}

	products, err := s.ProductRepo.FindByIds(businessId, productIDs)
	if err != nil {
		return nil, 0, err
	}

	var finalResults []response.ProductResponse
	for _, p := range products {
		res := helper.MapProductToResponse(p)

		if fromRedis, ok := autocompleteMap[p.Id]; ok && fromRedis.Image != nil && *fromRedis.Image != "" {
			res.Image = fromRedis.Image
		}

		finalResults = append(finalResults, res)
	}

	return finalResults, int64(len(finalResults)), nil
}

func (s *productService) SearchProductsRedisOnly(businessId int, search string, limit int) ([]response.ProductSearchResponse, error) {
	results, err := helper.GetProductAutocomplete(s.Redis, businessId, search, int64(limit))
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var searchResults []response.ProductSearchResponse
	for _, r := range results {
		searchResults = append(searchResults, response.ProductSearchResponse{
			Id:   r.Id,
			Name: r.Name,
		})
	}

	return searchResults, nil
}

func (s *productService) UpdateImage(id int, base64Image string) (response.ProductResponse, error) {
	product, err := s.ProductRepo.FindById(id)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("produk tidak ditemukan: %w", err)
	}

	oldImageURL := product.Image

	resultChan := make(chan struct {
		url string
		err error
	})

	go func() {
		url, err := helper.UploadBase64ToCloudinary(base64Image, "product")
		resultChan <- struct {
			url string
			err error
		}{url, err}
	}()

	result := <-resultChan
	if result.err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal upload gambar baru: %w", result.err)
	}
	product.Image = &result.url

	updatedProduct, err := s.ProductRepo.UpdateAll(&product)
	if err != nil {
		return response.ProductResponse{}, fmt.Errorf("gagal update produk: %w", err)
	}

	if oldImageURL != nil {
		go func(url string) {
			if publicID, err := helper.ExtractPublicIDFromURL(url); err == nil {
				_ = helper.DeleteFromCloudinary(publicID)
			}
		}(*oldImageURL)
	}

	return helper.MapProductToResponse(updatedProduct), nil
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

func (s *productService) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.ProductResponse, string, bool, error) {
	products, nextCursor, hasNext, err := s.ProductRepo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.ProductResponse
	for _, p := range products {
		result = append(result, helper.MapProductToResponse(p))
	}

	return result, nextCursor, hasNext, nil
}
