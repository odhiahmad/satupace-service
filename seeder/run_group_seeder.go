package seeders

import (
	"log"
	"time"

	"run-sync/entity"

	"run-sync/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRunGroups(db *gorm.DB) {
	groups := []entity.RunGroup{
		{
			Id:                uuid.New(),
			Name:              helper.StringPtr("Morning Running Group"),
			MinPace:           5.00,
			MaxPace:           6.00,
			PreferredDistance: 5,
			Latitude:          -6.2088,
			Longitude:         106.8456,
			ScheduledAt:       time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour).Add(6 * time.Hour),
			MaxMember:         20,
			IsWomenOnly:       false,
			Status:            "open",
			CreatedBy:         uuid.New(),
			CreatedAt:         time.Now(),
		},
		{
			Id:                uuid.New(),
			Name:              helper.StringPtr("Evening Running Group"),
			MinPace:           5.30,
			MaxPace:           6.30,
			PreferredDistance: 10,
			Latitude:          -6.2088,
			Longitude:         106.8456,
			ScheduledAt:       time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour).Add(18 * time.Hour),
			MaxMember:         25,
			IsWomenOnly:       false,
			Status:            "open",
			CreatedBy:         uuid.New(),
			CreatedAt:         time.Now(),
		},
		{
			Id:                uuid.New(),
			Name:              helper.StringPtr("Women Only Running Group"),
			MinPace:           5.30,
			MaxPace:           6.00,
			PreferredDistance: 5,
			Latitude:          -6.2088,
			Longitude:         106.8456,
			ScheduledAt:       time.Now().AddDate(0, 0, 2).Truncate(24 * time.Hour).Add(7 * time.Hour),
			MaxMember:         15,
			IsWomenOnly:       true,
			Status:            "open",
			CreatedBy:         uuid.New(),
			CreatedAt:         time.Now(),
		},
		{
			Id:                uuid.New(),
			Name:              helper.StringPtr("Half Marathon Training"),
			MinPace:           6.00,
			MaxPace:           7.00,
			PreferredDistance: 21,
			Latitude:          -6.2088,
			Longitude:         106.8456,
			ScheduledAt:       time.Now().AddDate(0, 0, 3).Truncate(24 * time.Hour).Add(6 * time.Hour),
			MaxMember:         30,
			IsWomenOnly:       false,
			Status:            "open",
			CreatedBy:         uuid.New(),
			CreatedAt:         time.Now(),
		},
	}

	for _, group := range groups {
		var existing entity.RunGroup
		err := db.Where("name = ?", helper.StringValue(group.Name)).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&group).Error; err != nil {
					log.Printf("Gagal menambahkan grup lari %s: %v", helper.StringValue(group.Name), err)
				} else {
					log.Printf("Grup lari %s berhasil ditambahkan", helper.StringValue(group.Name))
				}
			} else {
				log.Printf("Gagal memeriksa grup lari %s: %v", helper.StringValue(group.Name), err)
			}
		} else {
			log.Printf("Grup lari %s sudah ada, skip seeding", helper.StringValue(group.Name))
		}
	}
}
