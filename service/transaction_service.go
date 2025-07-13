package service

import (
	"errors"
	"fmt"
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
	Refund(itemReq request.TransactionRefundRequest) (*response.TransactionResponse, error)
	Cancel(itemReq request.TransactionRefundRequest) (*response.TransactionResponse, error)
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
		Status:     "unpaid",
		FinalPrice: res.FinalPrice,
		BasePrice:  res.BasePrice,
		SellPrice:  res.SellPrice,
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

func (s *transactionService) Payment(id int, req request.TransactionPaymentRequest) (*response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.AmountReceived == nil {
		return nil, errors.New("amountReceived harus diisi")
	}

	amount := *req.AmountReceived
	finalPrice := transaction.FinalPrice
	totalReceived := amount

	if transaction.AmountReceived != nil {
		totalReceived += *transaction.AmountReceived
	}

	var change float64
	var status string

	switch {
	case totalReceived >= finalPrice:
		status = "paid"
		change = totalReceived - finalPrice
	case totalReceived > 0 && totalReceived < finalPrice:
		status = "partial_paid"
		change = 0
	default:
		return nil, errors.New("jumlah pembayaran tidak valid")
	}

	now := time.Now().UTC()

	transaction.CustomerId = req.CustomerId
	transaction.PaymentMethodId = req.PaymentMethodId
	transaction.Rating = req.Rating
	transaction.Notes = req.Notes
	transaction.AmountReceived = &totalReceived
	transaction.Change = &change
	transaction.Status = status
	transaction.PaidAt = &now

	if status == "paid" {
		for _, item := range transaction.Items {
			qty := item.Quantity

			if item.ProductVariant != nil && item.ProductVariant.TrackStock {
				newStock := item.ProductVariant.Stock - qty
				if newStock < 0 {
					return nil, fmt.Errorf("stok produk varian tidak mencukupi: %s", item.ProductVariant.SKU)
				}
				if err := s.db.Model(&entity.ProductVariant{}).
					Where("id = ?", item.ProductVariant.Id).
					Update("stock", newStock).Error; err != nil {
					return nil, err
				}
			} else if item.Product != nil && item.Product.TrackStock {
				newStock := *item.Product.Stock - qty
				if newStock < 0 {
					return nil, fmt.Errorf("stok produk tidak mencukupi: %s", item.Product.Name)
				}
				if err := s.db.Model(&entity.Product{}).
					Where("id = ?", item.Product.Id).
					Update("stock", newStock).Error; err != nil {
					return nil, err
				}
			}
		}
	}

	savedTx, err := s.transactionRepo.Update(&transaction)
	if err != nil {
		return nil, err
	}

	return helper.MapTransactionResponse(savedTx), nil
}

func (s *transactionService) AddOrUpdateItem(transactionId int, itemReq request.TransactionItemCreate) (*response.TransactionResponse, error) {
	if itemReq.ProductId == nil && itemReq.BundleId == nil {
		return nil, errors.New("item harus memiliki product_id atau bundle_id")
	}

	var attrs []entity.TransactionItemAttribute
	for _, attr := range itemReq.Attributes {
		attrs = append(attrs, entity.TransactionItemAttribute{
			ProductAttributeId: attr.ProductAttributeId,
			AdditionalPrice:    attr.AdditionalPrice,
		})
	}

	item := entity.TransactionItem{
		ProductId:          itemReq.ProductId,
		BundleId:           itemReq.BundleId,
		ProductAttributeId: itemReq.ProductAttributeId,
		ProductVariantId:   itemReq.ProductVariantId,
		Quantity:           itemReq.Quantity,
		Attributes:         attrs,
	}

	if err := s.transactionRepo.AddOrReplaceItem(transactionId, item); err != nil {
		return nil, err
	}

	items, err := s.transactionRepo.FindItemsByTransactionId(transactionId)
	if err != nil {
		return nil, err
	}

	var allProductIds []int
	for _, itm := range items {
		if itm.ProductId != nil {
			allProductIds = append(allProductIds, *itm.ProductId)
		}
	}

	prepared, err := helper.PrepareTransactionItemsUpdate(helper.TransactionItemInputUpdate{
		DB:            s.db,
		Items:         helper.ToTransactionItemRequests(items),
		AllProductIds: allProductIds,
	})
	if err != nil {
		return nil, err
	}

	for _, updatedItem := range prepared.Items {
		if err := s.transactionRepo.UpdateItemFields(transactionId, updatedItem); err != nil {
			return nil, err
		}
	}

	update := &entity.Transaction{
		Id:         transactionId,
		FinalPrice: prepared.FinalPrice,
		SellPrice:  prepared.SellPrice,
		BasePrice:  prepared.SellPrice,
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

func (s *transactionService) FindById(id int) (*response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return helper.MapTransactionResponse(&transaction), nil
}

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

func (s *transactionService) Refund(itemReq request.TransactionRefundRequest) (*response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindById(itemReq.Id)
	if err != nil {
		return nil, err
	}

	if transaction.Status != "paid" && transaction.Status != "partial_paid" {
		return nil, errors.New("hanya transaksi dengan status 'paid' atau 'partial_paid' yang bisa direfund")
	}

	now := time.Now().UTC()

	transaction.Status = "refunded"
	transaction.IsRefunded = helper.BoolPtr(true)
	transaction.RefundReason = itemReq.Reason
	transaction.RefundedBy = &itemReq.UserId
	transaction.RefundedAt = &now

	savedTx, err := s.transactionRepo.Update(&transaction)
	if err != nil {
		return nil, err
	}

	return helper.MapTransactionResponse(savedTx), nil
}

func (s *transactionService) Cancel(itemReq request.TransactionRefundRequest) (*response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindById(itemReq.Id)
	if err != nil {
		return nil, err
	}

	if transaction.Status == "canceled" || transaction.Status == "refunded" {
		return nil, errors.New("transaksi sudah dibatalkan atau direfund sebelumnya")
	}

	if transaction.Status == "paid" {
		return nil, errors.New("transaksi yang sudah dibayar harus direfund, bukan dibatalkan")
	}

	now := time.Now().UTC()

	transaction.Status = "canceled"
	transaction.IsCanceled = helper.BoolPtr(true)
	transaction.CanceledReason = itemReq.Reason
	transaction.CanceledBy = &itemReq.UserId
	transaction.CanceledAt = &now

	savedTx, err := s.transactionRepo.Update(&transaction)
	if err != nil {
		return nil, err
	}

	return helper.MapTransactionResponse(savedTx), nil
}
