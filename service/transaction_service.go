package service

import (
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
	Create(req request.TransactionCreateRequest) (*entity.Transaction, error)
	Update(id int, req request.TransactionUpdateRequest) (*entity.Transaction, error)
	FindById(id int) (response.TransactionResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.TransactionResponse, int64, error)
	AddOrUpdateItem(transactionId int, item request.TransactionItemCreate) (*entity.Transaction, error)
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

func (s *transactionService) Create(req request.TransactionCreateRequest) (*entity.Transaction, error) {
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

	// Generate nomor bill
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
		Total:      res.Total,
		Discount:   res.TotalDiscount,
		Promo:      res.TotalPromo,
		CreatedAt:  time.Now(),
	}

	savedTx, err := s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, err
	}
	return savedTx, nil
}

// UPDATE
func (s *transactionService) Update(id int, req request.TransactionUpdateRequest) (*entity.Transaction, error) {
	var allProductIds []int
	for _, item := range req.Items {
		productId := helper.IntOrDefault(item.ProductId, 0)
		if productId != 0 {
			allProductIds = append(allProductIds, productId)
		}
	}

	res, err := helper.PrepareTransactionItemsUpdate(helper.TransactionItemInputUpdate{
		DB:            s.db,
		Items:         req.Items,
		AllProductIds: allProductIds,
	})
	if err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		Id:              id,
		CustomerId:      req.CustomerId,
		PaymentMethodId: req.PaymentMethodId,
		BillNumber:      req.BillNumber,
		Items:           res.Items,
		Total:           res.Total,
		Discount:        res.TotalDiscount,
		Promo:           res.TotalPromo,
		Status:          req.Status,
		Rating:          req.Rating,
		Notes:           req.Notes,
		AmountReceived:  req.AmountReceived,
		Change:          req.Change,
	}

	savedTx, err := s.transactionRepo.Update(transaction)
	if err != nil {
		return nil, err
	}
	return savedTx, nil
}

func (s *transactionService) AddOrUpdateItem(transactionId int, itemReq request.TransactionItemCreate) (*entity.Transaction, error) {
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
		Price:              itemReq.Price,
		Discount:           itemReq.Discount,
		Promo:              itemReq.Promo,
		Rating:             itemReq.Rating,
		Attributes:         attrs,
	}

	_, err := s.transactionRepo.AddOrUpdateItem(transactionId, item)
	if err != nil {
		return nil, err
	}

	items, err := s.transactionRepo.FindItemsByTransactionId(transactionId)
	if err != nil {
		return nil, err
	}

	res, err := helper.CalculateTransactionTotals(s.db, items)
	if err != nil {
		return nil, err
	}

	update := &entity.Transaction{
		Id:       transactionId,
		Total:    res.Total,
		Discount: res.TotalDiscount,
		Promo:    res.TotalPromo,
	}

	return s.transactionRepo.Update(update)
}

// FINDBYID
func (s *transactionService) FindById(id int) (response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindById(id)
	if err != nil {
		return response.TransactionResponse{}, err
	}
	return ToTransactionResponse(transaction), nil
}

// FINDWITHPAGINATION
func (s *transactionService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.TransactionResponse, int64, error) {
	transactions, total, err := s.transactionRepo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []response.TransactionResponse
	for _, trx := range transactions {
		responses = append(responses, ToTransactionResponse(trx))
	}

	return responses, total, nil
}

func ToTransactionResponse(trx entity.Transaction) response.TransactionResponse {
	var itemResponses []response.TransactionItemResponse

	for _, item := range trx.Items {
		var attrResponses []response.TransactionItemAttributeResponse
		for _, attr := range item.Attributes {
			attrResponses = append(attrResponses, response.TransactionItemAttributeResponse{
				Id:                 attr.Id,
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		itemResponses = append(itemResponses, response.TransactionItemResponse{
			Id:                 item.Id,
			ProductId:          item.ProductId,
			BundleId:           item.BundleId,
			ProductAttributeId: item.ProductAttributeId,
			ProductVariantId:   item.ProductVariantId,
			Quantity:           item.Quantity,
			Price:              item.Price,
			Discount:           item.Discount,
			Promo:              item.Promo,
			Rating:             item.Rating,
			Attributes:         attrResponses,
		})
	}

	return response.TransactionResponse{
		Id:              trx.Id,
		BusinessId:      trx.BusinessId,
		CustomerId:      trx.CustomerId,
		PaymentMethodId: trx.PaymentMethodId,
		BillNumber:      trx.BillNumber,
		Items:           itemResponses,
		Total:           trx.Total,
		Discount:        trx.Discount,
		Promo:           trx.Promo,
		Status:          trx.Status,
		Rating:          trx.Rating,
		Notes:           trx.Notes,
		AmountReceived:  trx.AmountReceived,
		Change:          trx.Change,
	}
}
