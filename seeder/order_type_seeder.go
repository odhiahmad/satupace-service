package seeders

import (
	"log"
	"time"

	"loka-kasir/entity"

	"gorm.io/gorm"
)

func SeedOrderTypes(db *gorm.DB) {
	orderTypes := []entity.OrderType{
		{
			Code:      "DIN",
			Name:      "dine_in",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Code:      "TKA",
			Name:      "take_away",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Code:      "DLV",
			Name:      "delivery",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, orderType := range orderTypes {
		var existing entity.OrderType
		err := db.Where("name = ?", orderType.Name).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&orderType).Error; err != nil {
					log.Printf("❌ Gagal menambahkan order type %s: %v", orderType.Name, err)
				} else {
					log.Printf("✅ Order type %s berhasil ditambahkan", orderType.Name)
				}
			} else {
				log.Printf("⚠️ Gagal memeriksa order type %s: %v", orderType.Name, err)
			}
		} else {
			log.Printf("⏭️ Order type %s sudah ada, skip seeding", orderType.Name)
		}
	}
}
