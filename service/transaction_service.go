package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type TransactionService interface {
	Create(req request.TransactionCreateRequest) error
	Update(id int, req request.TransactionUpdateRequest) error
	FindById(id int) (response.TransactionResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.TransactionResponse, int64, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	validate        *validator.Validate
}

func NewTransactionService(repo repository.TransactionRepository, validate *validator.Validate) TransactionService {
	return &transactionService{
		transactionRepo: repo,
		validate:        validator.New(),
	}
}

func (s *transactionService) Create(req request.TransactionCreateRequest) error {
	var items []entity.TransactionItem

	for _, item := range req.Items {
		var attrs []entity.TransactionItemAttribute
		for _, attr := range item.Attributes {
			attrs = append(attrs, entity.TransactionItemAttribute{
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		items = append(items, entity.TransactionItem{
			ProductId:          item.ProductId,
			BundleId:           item.BundleId,
			ProductAttributeId: item.ProductAttributeId,
			ProductVariantId:   item.ProductVariantId,
			Quantity:           item.Quantity,
			UnitPrice:          item.UnitPrice,
			Price:              item.Price,
			Discount:           item.Discount,
			Promo:              item.Promo,
			Rating:             item.Rating,
			Attributes:         attrs,
		})
	}

	transaction := &entity.Transaction{
		BusinessId:      req.BusinessId,
		CustomerId:      req.CustomerId,
		PaymentMethodId: req.PaymentMethodId,
		BillNumber:      req.BillNumber,
		Items:           items,
		Total:           req.Total,
		Discount:        req.Discount,
		Promo:           req.Promo,
		Status:          req.Status,
		Rating:          req.Rating,
		Notes:           req.Notes,
		AmountReceived:  req.AmountReceived,
		Change:          req.Change,
	}

	return s.transactionRepo.Create(transaction)
}

// UPDATE
func (s *transactionService) Update(id int, req request.TransactionUpdateRequest) error {
	var items []entity.TransactionItem

	for _, item := range req.Items {
		var attrs []entity.TransactionItemAttribute
		for _, attr := range item.Attributes {
			attrs = append(attrs, entity.TransactionItemAttribute{
				Id:                 attr.Id,
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		items = append(items, entity.TransactionItem{
			Id:                 item.Id,
			ProductId:          item.ProductId,
			BundleId:           item.BundleId,
			ProductAttributeId: item.ProductAttributeId,
			ProductVariantId:   item.ProductVariantId,
			Quantity:           item.Quantity,
			UnitPrice:          item.UnitPrice,
			Price:              item.Price,
			Discount:           item.Discount,
			Promo:              item.Promo,
			Rating:             item.Rating,
			Attributes:         attrs,
		})
	}

	transaction := &entity.Transaction{
		Id:              id,
		CustomerId:      req.CustomerId,
		PaymentMethodId: req.PaymentMethodId,
		BillNumber:      req.BillNumber,
		Items:           items,
		Total:           req.Total,
		Discount:        req.Discount,
		Promo:           req.Promo,
		Status:          req.Status,
		Rating:          req.Rating,
		Notes:           req.Notes,
		AmountReceived:  req.AmountReceived,
		Change:          req.Change,
	}

	return s.transactionRepo.Update(transaction)
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
			UnitPrice:          item.UnitPrice,
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
		CreatedAt:       trx.CreatedAt,
		UpdatedAt:       trx.UpdatedAt,
	}
}
