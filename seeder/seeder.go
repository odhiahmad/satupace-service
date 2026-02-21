package seeders

import "gorm.io/gorm"

func RunAll(db *gorm.DB) {
	RunAllSeeders(db)
}
