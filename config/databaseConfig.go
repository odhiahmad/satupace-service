package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setup database
func SetupDatabaseConnection() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		panic("Gagal load file env")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Shanghai")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Gagal membuat koneksi ke database")
	}

	db.AutoMigrate(
		&entity.User{},
		&entity.Business{},
		&entity.BusinessType{},
		&entity.Role{},
		&entity.Customer{},
		&entity.Product{},
		&entity.ProductSize{},
		&entity.ProductUnit{},
		&entity.MenuAttribute{},
		&entity.Transaction{},
		&entity.TransactionDetail{},
		&entity.PaymentMethod{})

	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()

	if err != nil {
		panic("Gagal menutup koneksi dari database")
	}

	dbSQL.Close()
}
