package seeders

import (
	"log"
	"time"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

func SeedPaymentMethods(db *gorm.DB) {
	paymentMethods := []entity.PaymentMethod{
		{Name: "Cash", Code: "CASH"},
		{Name: "Debit Card", Code: "DEBIT"},
		{Name: "Credit Card", Code: "CREDIT"},
		{Name: "QRIS", Code: "QRIS"},
		{Name: "E-Wallet", Code: "EWALLET"},
	}

	for _, pm := range paymentMethods {
		var existing entity.PaymentMethod
		err := db.Where("code = ?", pm.Code).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				pm.CreatedAt = time.Now()
				pm.UpdatedAt = time.Now()
				if err := db.Create(&pm).Error; err != nil {
					log.Printf("Gagal menambahkan payment method %s: %v", pm.Name, err)
				} else {
					log.Printf("Payment method %s berhasil ditambahkan", pm.Name)
				}
			} else {
				log.Printf("Gagal memeriksa payment method %s: %v", pm.Name, err)
			}
		} else {
			log.Printf("Payment method %s sudah ada, skip seeding", pm.Name)
		}
	}
}
