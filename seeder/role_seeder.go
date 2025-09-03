package seeders

import (
	"log"
	"time"

	"loka-kasir/entity"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []entity.Role{
		{Name: "Owner"},
		{Name: "Manager"},
		{Name: "Kasir"},
		{Name: "Barista"},
		{Name: "Koki"},
		{Name: "Pelayan"},
		{Name: "Gudang"},
		{Name: "Kurir"},
		{Name: "Supervisor"},
		{Name: "Admin"},
	}

	for _, role := range roles {
		var existing entity.Role
		err := db.Where("name = ?", role.Name).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				role.CreatedAt = time.Now()
				role.UpdatedAt = time.Now()
				if err := db.Create(&role).Error; err != nil {
					log.Printf("Gagal menambahkan role %s: %v", role.Name, err)
				} else {
					log.Printf("Role %s berhasil ditambahkan", role.Name)
				}
			} else {
				log.Printf("Gagal memeriksa role %s: %v", role.Name, err)
			}
		} else {
			log.Printf("Role %s sudah ada, skip seeding", role.Name)
		}
	}
}
