package repository

import (
	"errors"
	"fmt"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) (*entity.Transaction, error)
	Update(transaction *entity.Transaction) (*entity.Transaction, error)
	FindById(id int) (entity.Transaction, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Transaction, int64, error)
	AddOrReplaceItem(transactionId int, item entity.TransactionItem) error
	FindItemsByTransactionId(transactionId int) ([]entity.TransactionItem, error)
	UpdateTotals(transaction *entity.Transaction) (*entity.Transaction, error)
	UpdateItemFields(transactionId int, item entity.TransactionItem) error
}

type transactionConnection struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionConnection{db}
}

// CREATE
func (conn *transactionConnection) Create(transaction *entity.Transaction) (*entity.Transaction, error) {
	err := conn.db.Transaction(func(tx *gorm.DB) error {
		items := transaction.Items
		transaction.Items = nil

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		transaction.Items = items

		fmt.Println("Jumlah item yang akan disimpan:", len(transaction.Items))
		for i, item := range transaction.Items {
			fmt.Printf("Item: %+v\n", item)
			transaction.Items[i].Id = 0
			transaction.Items[i].TransactionId = transaction.Id

			if err := tx.Create(&transaction.Items[i]).Error; err != nil {
				return err
			}

			for j := range transaction.Items[i].Attributes {
				transaction.Items[i].Attributes[j].Id = 0
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
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Preload("Items.Product.ProductPromos").
		Preload("Items.Product.ProductPromos.Promo").
		Preload("Items.Product.ProductPromos.Promo.RequiredProducts").
		First(&result, transaction.Id).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

// UPDATE
func (conn *transactionConnection) Update(transaction *entity.Transaction) (*entity.Transaction, error) {
	err := conn.db.Transaction(func(tx *gorm.DB) error {
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

	if err != nil {
		return nil, err
	}

	// ✅ Reload data transaksi dengan relasi lengkap setelah berhasil disimpan
	var result entity.Transaction
	if err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Preload("Items.Product.ProductPromos").
		Preload("Items.Product.ProductPromos.Promo").
		Preload("Items.Product.ProductPromos.Promo.RequiredProducts").
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

	// preload relasi jika diperlukan
	if err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Preload("Items.Product.ProductPromos").
		Preload("Items.Product.ProductPromos.Promo").
		Preload("Items.Product.ProductPromos.Promo.RequiredProducts").
		First(&existing, existing.Id).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (conn *transactionConnection) UpdateItemFields(transactionId int, item entity.TransactionItem) error {
	if item.ProductId == nil {
		return fmt.Errorf("product_id kosong pada item")
	}

	return conn.db.Model(&entity.TransactionItem{}).
		Where("transaction_id = ? AND product_id = ?", transactionId, *item.ProductId).
		Updates(map[string]interface{}{
			"quantity":   item.Quantity,
			"unit_price": item.UnitPrice,
			"total":      item.Total,
			"discount":   item.Discount,
			"promo":      item.Promo,
			"tax":        item.Tax,
			"rating":     item.Rating,
		}).Error
}

// FINDBYID
func (conn *transactionConnection) FindById(id int) (entity.Transaction, error) {
	var transaction entity.Transaction
	err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Preload("Items.Product.ProductPromos").
		Preload("Items.Product.ProductPromos.Promo").
		Preload("Items.Product.ProductPromos.Promo.RequiredProducts").
		Where("id = ?", id).
		First(&transaction).Error
	return transaction, err
}

// FINDWITHPAGINATION
func (conn *transactionConnection) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64

	// Base query dengan preload relasi dan sorting
	baseQuery := conn.db.Model(&entity.Transaction{}).
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Preload("Items.Product.ProductPromos").
		Preload("Items.Product.ProductPromos.Promo").
		Preload("Items.Product.ProductPromos.Promo.RequiredProducts").
		Where("business_id = ?", businessId)

	// Hitung total data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Siapkan paginator
	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	// Jalankan paginasi
	_, _, err := p.Paginate(baseQuery, &transactions)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (conn *transactionConnection) AddOrReplaceItem(transactionId int, item entity.TransactionItem) error {
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
			return err // error selain record not found
		}

		// Item sudah ada
		if item.Quantity <= 0 {
			// Hapus jika quantity <= 0
			if err := tx.Where("transaction_item_id = ?", existing.Id).
				Delete(&entity.TransactionItemAttribute{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
			return nil
		}

		// Update isi item
		existing.Quantity = item.Quantity
		existing.UnitPrice = item.UnitPrice
		existing.Total = item.Total
		existing.Discount = item.Discount
		existing.Promo = item.Promo
		existing.Tax = item.Tax
		existing.Rating = item.Rating

		// Hapus atribut lama
		if err := tx.Where("transaction_item_id = ?", existing.Id).
			Delete(&entity.TransactionItemAttribute{}).Error; err != nil {
			return err
		}

		// Tambah atribut baru
		for i := range item.Attributes {
			item.Attributes[i].TransactionItemId = existing.Id
			if err := tx.Create(&item.Attributes[i]).Error; err != nil {
				return err
			}
		}

		// Simpan item yang diupdate
		return tx.Save(&existing).Error
	})
}

func (conn *transactionConnection) FindItemsByTransactionId(transactionId int) ([]entity.TransactionItem, error) {
	var items []entity.TransactionItem

	err := conn.db.
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Preload("Items.Bundle").
		Preload("Items.ProductAttribute").
		Preload("Items.ProductVariant").
		Preload("Items.Attributes.ProductAttribute").
		Preload("Items.Product.Tax").
		Preload("Items.Product.Discount").
		Preload("Items.Product.Unit").
		Preload("Items.Product.ProductPromos").
		Preload("Items.Product.ProductPromos.Promo").
		Preload("Items.Product.ProductPromos.Promo.RequiredProducts").
		Where("transaction_id = ?", transactionId).
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}
