package seeders

import (
	"log"
	"time"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

func SeedBusinessTypes(db *gorm.DB) {
	businessTypes := []entity.BusinessType{
		{Name: "Restoran"},
		{Name: "Kafe"},
		{Name: "Toko Kelontong"},
		{Name: "Minimarket"},
		{Name: "Apotek"},
		{Name: "Bakery"},
		{Name: "Barbershop"},
		{Name: "Bengkel"},
		{Name: "Butik"},
		{Name: "Toko Elektronik"},
	}

	for _, bt := range businessTypes {
		var existing entity.BusinessType
		err := db.Where("name = ?", bt.Name).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				bt.CreatedAt = time.Now()
				bt.UpdatedAt = time.Now()
				if err := db.Create(&bt).Error; err != nil {
					log.Printf("Gagal menambahkan business type %s: %v", bt.Name, err)
				} else {
					log.Printf("Business type %s berhasil ditambahkan", bt.Name)
				}
			} else {
				log.Printf("Gagal memeriksa business type %s: %v", bt.Name, err)
			}
		} else {
			log.Printf("Business type %s sudah ada, skip seeding", bt.Name)
		}
	}
}
