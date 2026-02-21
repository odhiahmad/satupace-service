package seeders

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) {
	log.Println("================================")
	log.Println("Running Seeders")
	log.Println("================================")

	SeedRunGroups(db)

	log.Println("================================")
	log.Println("Seeding Finished")
	log.Println("================================")
}
