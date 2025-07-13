package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
	"gorm.io/gorm"
)

type TransactionService interface {
	Create(req request.TransactionCreateRequest) (*response.TransactionResponse, error)
	Payment(id int, req request.TransactionPaymentRequest) (*response.TransactionResponse, error)
	FindById(id int) (*response.TransactionResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]*response.TransactionResponse, int64, error)
	AddOrUpdateItem(transactionId int, item request.TransactionItemCreate) (*response.TransactionResponse, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	validate        *validator.Validate
	db              *gorm.DB
}

func NewTransactionService(db *gorm.DB, repo repository.TransactionRepository, validate *validator.Validate) TransactionService {
	return &transactionService{
		db:              db,
		transactionRepo: repo,
		validate:        validator.New(),
	}
}

func (s *transactionService) Create(req request.TransactionCreateRequest) (*response.TransactionResponse, error) {
	var allProductIds []int
	for _, item := range req.Items {
		productId := helper.IntOrDefault(item.ProductId, 0)
		if productId != 0 {
			allProductIds = append(allProductIds, productId)
		}
	}

	res, err := helper.PrepareTransactionItemsCreate(helper.TransactionItemInput{
		DB:            s.db,
		Items:         req.Items,
		AllProductIds: allProductIds,
	})

	if err != nil {
		return nil, err
	}

	billNumber, err := helper.GenerateBillNumber(s.db)
	if err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		BusinessId: req.BusinessId,
		CustomerId: req.CustomerId,
		BillNumber: billNumber,
		Items:      res.Items,
		Status:     "cart",
		FinalPrice: res.FinalPrice,
		BasePrice:  res.BasePrice,
		Discount:   res.TotalDiscount,
		Tax:        res.TotalTax,
		CreatedAt:  time.Now(),
	}

	savedTx, err := s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, err
	}

	return helper.MapTransactionResponse(savedTx), nil
}

// UPDATE
func (s *transactionService) Payment(id int, req request.TransactionPaymentRequest) (*response.TransactionResponse, error) {
	// Ambil transaksi dari database
	transaction, err := s.transactionRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	// Hitung kembalian jika AmountReceived tersedia
	var change float64
	if req.AmountReceived != nil {
		change = *req.AmountReceived - transaction.FinalPrice
	}

	// Update field transaksi
	transaction.CustomerId = req.CustomerId
	transaction.PaymentMethodId = req.PaymentMethodId
	transaction.Rating = req.Rating
	transaction.Notes = req.Notes
	transaction.AmountReceived = req.AmountReceived
	transaction.Change = &change
	transaction.Status = "paid"
	transaction.PaidAt = time.Now().UTC()

	// Simpan perubahan ke database
	savedTx, err := s.transactionRepo.Update(&transaction)
	if err != nil {
		return nil, err
	}

	// Kembalikan response
	return helper.MapTransactionResponse(savedTx), nil
}

func (s *transactionService) AddOrUpdateItem(transactionId int, itemReq request.TransactionItemCreate) (*response.TransactionResponse, error) {
	if itemReq.ProductId == nil && itemReq.BundleId == nil {
		return nil, errors.New("item harus memiliki product_id atau bundle_id")
	}

	// Konversi atribut
	var attrs []entity.TransactionItemAttribute
	for _, attr := range itemReq.Attributes {
		attrs = append(attrs, entity.TransactionItemAttribute{
			ProductAttributeId: attr.ProductAttributeId,
			AdditionalPrice:    attr.AdditionalPrice,
		})
	}

	// Buat entity
	item := entity.TransactionItem{
		ProductId:          itemReq.ProductId,
		BundleId:           itemReq.BundleId,
		ProductAttributeId: itemReq.ProductAttributeId,
		ProductVariantId:   itemReq.ProductVariantId,
		Quantity:           itemReq.Quantity,
		Attributes:         attrs,
	}

	// Tambah atau update item di DB
	if err := s.transactionRepo.AddOrReplaceItem(transactionId, item); err != nil {
		return nil, err
	}

	// Ambil ulang semua item setelah update
	items, err := s.transactionRepo.FindItemsByTransactionId(transactionId)
	if err != nil {
		return nil, err
	}

	// Ambil semua product_id
	var allProductIds []int
	for _, itm := range items {
		if itm.ProductId != nil {
			allProductIds = append(allProductIds, *itm.ProductId)
		}
	}

	// Hitung ulang untuk promo, diskon, dan harga final
	prepared, err := helper.PrepareTransactionItemsUpdate(helper.TransactionItemInputUpdate{
		DB:            s.db,
		Items:         helper.ToTransactionItemRequests(items),
		AllProductIds: allProductIds,
	})
	if err != nil {
		return nil, err
	}

	// Update masing-masing item
	for _, updatedItem := range prepared.Items {
		if err := s.transactionRepo.UpdateItemFields(transactionId, updatedItem); err != nil {
			return nil, err
		}
	}

	// Update total transaksi saja
	update := &entity.Transaction{
		Id:         transactionId,
		FinalPrice: prepared.FinalPrice,
		BasePrice:  prepared.BasePrice,
		Discount:   prepared.TotalDiscount,
		Promo:      prepared.TotalPromo,
		Tax:        prepared.TotalTax,
	}
	savedTx, err := s.transactionRepo.UpdateTotals(update)
	if err != nil {
		return nil, err
	}

	return helper.MapTransactionResponse(savedTx), nil
}

// FINDBYID
func (s *transactionService) FindById(id int) (*response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return helper.MapTransactionResponse(&transaction), nil
}

// FINDWITHPAGINATION
func (s *transactionService) FindWithPagination(businessId int, pagination request.Pagination) ([]*response.TransactionResponse, int64, error) {
	transactions, total, err := s.transactionRepo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []*response.TransactionResponse
	for _, trx := range transactions {
		responses = append(responses, helper.MapTransactionResponse(&trx)) // kalau trx adalah entity.Transaction
	}

	return responses, total, nil
}
