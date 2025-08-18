package seeders

import "gorm.io/gorm"

func RunAll(db *gorm.DB) {
	SeedOrderTypes(db)
	SeedBusinessTypes(db)
	SeedRoles(db)
	SeedPaymentMethods(db)
}
