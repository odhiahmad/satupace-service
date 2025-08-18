package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) (*entity.Transaction, error)
	Update(transaction *entity.Transaction) (*entity.Transaction, error)
	FindById(id uuid.UUID) (entity.Transaction, error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Transaction, int64, error)
	AddOrReplaceItem(transactionId uuid.UUID, item entity.TransactionItem) error
	FindItemsByTransactionId(transactionId uuid.UUID) ([]entity.TransactionItem, error)
	UpdateTotals(transaction *entity.Transaction) (*entity.Transaction, error)
	UpdateItemFields(transactionId uuid.UUID, item entity.TransactionItem) error
}

type transactionConnection struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionConnection{db}
}

func (conn *transactionConnection) Create(transaction *entity.Transaction) (*entity.Transaction, error) {
	err := conn.db.Transaction(func(tx *gorm.DB) error {
		items := transaction.Items
		transaction.Items = nil

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].Id = uuid.Nil
			items[i].TransactionId = transaction.Id
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		var allAttributes []entity.TransactionItemAttribute
		for _, item := range items {
			for j := range item.Attributes {
				attr := &item.Attributes[j]
				attr.Id = uuid.Nil
				attr.TransactionItemId = item.Id
				allAttributes = append(allAttributes, *attr)
			}
		}

		if len(allAttributes) > 0 {
			if err := tx.Create(&allAttributes).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var result entity.Transaction
	if err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Preload("Cashier").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		First(&result, transaction.Id).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (conn *transactionConnection) Update(transaction *entity.Transaction) (*entity.Transaction, error) {
	err := conn.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Transaction{}).
			Where("id = ?", transaction.Id).
			Updates(map[string]interface{}{
				"cashier_id":        transaction.CashierId,
				"customer_id":       transaction.CustomerId,
				"payment_method_id": transaction.PaymentMethodId,
				"rating":            transaction.Rating,
				"amount_received":   transaction.AmountReceived,
				"change":            transaction.Change,
				"status":            transaction.Status,
				"paid_at":           transaction.PaidAt,
			}).Error; err != nil {
			return err
		}

		if err := tx.Where("transaction_item_id IN (?)",
			tx.Model(&entity.TransactionItem{}).Select("id").Where("transaction_id = ?", transaction.Id),
		).Delete(&entity.TransactionItemAttribute{}).Error; err != nil {
			return err
		}

		if err := tx.Where("transaction_id = ?", transaction.Id).Delete(&entity.TransactionItem{}).Error; err != nil {
			return err
		}

		for i := range transaction.Items {
			transaction.Items[i].Id = uuid.Nil
			transaction.Items[i].TransactionId = transaction.Id

			if err := tx.Create(&transaction.Items[i]).Error; err != nil {
				return err
			}

			for j := range transaction.Items[i].Attributes {
				transaction.Items[i].Attributes[j].Id = uuid.Nil
				transaction.Items[i].Attributes[j].TransactionItemId = transaction.Items[i].Id
				if err := tx.Create(&transaction.Items[i].Attributes[j]).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var result entity.Transaction
	if err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Cashier").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		First(&result, transaction.Id).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (conn *transactionConnection) UpdateTotals(transaction *entity.Transaction) (*entity.Transaction, error) {
	var existing entity.Transaction
	if err := conn.db.First(&existing, transaction.Id).Error; err != nil {
		return nil, err
	}

	existing.FinalPrice = transaction.FinalPrice
	existing.BasePrice = transaction.BasePrice
	existing.Discount = transaction.Discount
	existing.Promo = transaction.Promo

	if err := conn.db.Save(&existing).Error; err != nil {
		return nil, err
	}

	if err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Cashier").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		First(&existing, existing.Id).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (conn *transactionConnection) UpdateItemFields(transactionId uuid.UUID, item entity.TransactionItem) error {
	if item.ProductId == nil {
		return fmt.Errorf("product_id kosong pada item")
	}

	return conn.db.Model(&entity.TransactionItem{}).
		Where("transaction_id = ? AND product_id = ?", transactionId, *item.ProductId).
		Updates(map[string]interface{}{
			"quantity":   item.Quantity,
			"base_price": item.BasePrice,
			"sell_price": item.SellPrice,
			"total":      item.Total,
			"discount":   item.Discount,
			"promo":      item.Promo,
			"tax":        item.Tax,
			"rating":     item.Rating,
		}).Error
}

func (conn *transactionConnection) FindById(id uuid.UUID) (entity.Transaction, error) {
	var transaction entity.Transaction
	err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Cashier").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Where("id = ?", id).
		First(&transaction).Error
	return transaction, err
}

func (conn *transactionConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64

	baseQuery := conn.db.Model(&entity.Transaction{}).
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Cashier").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Where("business_id = ?", businessId)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &transactions)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (conn *transactionConnection) AddOrReplaceItem(transactionId uuid.UUID, item entity.TransactionItem) error {
	return conn.db.Transaction(func(tx *gorm.DB) error {
		var existing entity.TransactionItem

		query := tx.Where("transaction_id = ?", transactionId)

		if item.ProductId != nil {
			query = query.Where("product_id = ?", item.ProductId)
		} else {
			query = query.Where("product_id IS NULL")
		}

		if item.BundleId != nil {
			query = query.Where("bundle_id = ?", item.BundleId)
		} else {
			query = query.Where("bundle_id IS NULL")
		}

		if item.ProductVariantId != nil {
			query = query.Where("product_variant_id = ?", item.ProductVariantId)
		} else {
			query = query.Where("product_variant_id IS NULL")
		}

		if item.ProductAttributeId != nil {
			query = query.Where("product_attribute_id = ?", item.ProductAttributeId)
		} else {
			query = query.Where("product_attribute_id IS NULL")
		}

		err := query.First(&existing).Error

		// Item belum ada → insert jika quantity > 0
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if item.Quantity > 0 {
				item.TransactionId = transactionId
				if err := tx.Create(&item).Error; err != nil {
					return err
				}
				for i := range item.Attributes {
					item.Attributes[i].TransactionItemId = item.Id
					if err := tx.Create(&item.Attributes[i]).Error; err != nil {
						return err
					}
				}
			}
			return nil
		}

		if err != nil {
			return err
		}

		if item.Quantity <= 0 {
			if err := tx.Where("transaction_item_id = ?", existing.Id).
				Delete(&entity.TransactionItemAttribute{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
			return nil
		}

		existing.Quantity = item.Quantity
		existing.BasePrice = item.BasePrice
		existing.SellPrice = item.SellPrice
		existing.Total = item.Total
		existing.Discount = item.Discount
		existing.Promo = item.Promo
		existing.Tax = item.Tax
		existing.Rating = item.Rating

		if err := tx.Where("transaction_item_id = ?", existing.Id).
			Delete(&entity.TransactionItemAttribute{}).Error; err != nil {
			return err
		}

		for i := range item.Attributes {
			item.Attributes[i].TransactionItemId = existing.Id
			if err := tx.Create(&item.Attributes[i]).Error; err != nil {
				return err
			}
		}

		return tx.Save(&existing).Error
	})
}

func (conn *transactionConnection) FindItemsByTransactionId(transactionId uuid.UUID) ([]entity.TransactionItem, error) {
	var items []entity.TransactionItem

	err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Cashier").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Where("transaction_id = ?", transactionId).
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}
