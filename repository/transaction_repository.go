package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) error
	Update(transaction *entity.Transaction) error
	FindById(id int) (entity.Transaction, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Transaction, int64, error)
	AddOrUpdateItem(transactionId int, item entity.TransactionItem) error
	FindItemsByTransactionId(transactionId int) ([]entity.TransactionItem, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

// CREATE
func (r *transactionRepository) Create(transaction *entity.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		for i := range transaction.Items {
			transaction.Items[i].TransactionId = transaction.Id
			if err := tx.Create(&transaction.Items[i]).Error; err != nil {
				return err
			}

			for j := range transaction.Items[i].Attributes {
				transaction.Items[i].Attributes[j].TransactionItemId = transaction.Items[i].Id
				if err := tx.Create(&transaction.Items[i].Attributes[j]).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// UPDATE
func (r *transactionRepository) Update(transaction *entity.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(transaction).Error; err != nil {
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
			transaction.Items[i].TransactionId = transaction.Id
			if err := tx.Create(&transaction.Items[i]).Error; err != nil {
				return err
			}

			for j := range transaction.Items[i].Attributes {
				transaction.Items[i].Attributes[j].TransactionItemId = transaction.Items[i].Id
				if err := tx.Create(&transaction.Items[i].Attributes[j]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// FINDBYID
func (r *transactionRepository) FindById(id int) (entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Where("id = ?", id).
		First(&transaction).Error
	return transaction, err
}

// FINDWITHPAGINATION
func (r *transactionRepository) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64

	// Base query dengan preload relasi dan sorting
	baseQuery := r.db.Model(&entity.Transaction{}).
		Where("business_id = ?", businessId).
		Preload("Customer").
		Preload("PaymentMethod").
		Order("created_at desc")

	// Hitung total data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination)

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &transactions)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *transactionRepository) AddOrUpdateItem(transactionId int, item entity.TransactionItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing entity.TransactionItem
		err := tx.Where("transaction_id = ? AND product_id = ?", transactionId, item.ProductId).First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Insert baru
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
		} else {
			// Update kuantitas dan subtotal
			existing.Quantity += item.Quantity
			existing.Price = item.Price // atau logic subtotal baru?
			if err := tx.Save(&existing).Error; err != nil {
				return err
			}

			// Hapus atribut lama
			if err := tx.Where("transaction_item_id = ?", existing.Id).Delete(&entity.TransactionItemAttribute{}).Error; err != nil {
				return err
			}

			// Tambahkan atribut baru
			for i := range item.Attributes {
				item.Attributes[i].TransactionItemId = existing.Id
				if err := tx.Create(&item.Attributes[i]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *transactionRepository) FindItemsByTransactionId(transactionId int) ([]entity.TransactionItem, error) {
	var items []entity.TransactionItem

	err := r.db.
		Preload("Attributes").
		Where("transaction_id = ?", transactionId).
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}
