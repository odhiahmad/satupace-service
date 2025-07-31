package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/vandyahmad24/golang-wilayah-indonesia/wilayah"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabaseConnection() *gorm.DB {
	ginMode := os.Getenv("GIN_MODE")

	if ginMode != "release" {
		if err := godotenv.Load(); err != nil {
			log.Println("Peringatan: .env file tidak ditemukan, menggunakan env bawaan sistem.")
		}
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" || dbPort == "" {
		log.Fatal("❌ Environment variable database tidak lengkap")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}

	if ginMode != "release" {
		if err := db.AutoMigrate(
			&entity.UserBusiness{},
			&entity.BusinessType{},
			&entity.Business{},
			&entity.Membership{},
			&entity.Category{},
			&entity.Unit{},
			&entity.Tax{},
			&entity.Discount{},
			&entity.Product{},
			&entity.ProductVariant{},
			&entity.ProductAttribute{},
			&entity.Bundle{},
			&entity.BundleItem{},
			&entity.PaymentMethod{},
			&entity.Customer{},
			&entity.Transaction{},
			&entity.TransactionItem{},
			&entity.TransactionItemAttribute{},
			&entity.Brand{},
		); err != nil {
			log.Fatalf("❌ AutoMigrate gagal: %v", err)
		}
	} else {
		log.Println("ℹ️ Production mode terdeteksi, AutoMigrate dilewati.")
	}

	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan koneksi SQL: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Gagal menutup koneksi database: %v", err)
	}
}

func SetupWhatsAppGORM() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME_WHATSAPP")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke DB WhatsApp: %v", err)
	}

	return db
}

func SetupWilayahDatabase() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi SQL: %v", err)
	}
	defer db.Close()

	wilayah.RunMigration(db)
	wilayah.Seed(db, "data")

}
