package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabaseConnection menginisialisasi koneksi ke database PostgreSQL.
func SetupDatabaseConnection() *gorm.DB {
	appEnv := os.Getenv("APP_ENV") // development | release | staging

	if appEnv != "release" {
		if err := godotenv.Load(); err != nil {
			log.Println("Peringatan: .env file tidak ditemukan, menggunakan env bawaan sistem.")
		}
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Validasi env
	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" || dbPort == "" {
		log.Fatal("❌ Environment variable database tidak lengkap")
	}

	// Format DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}

	// ENUM setup, bisa kamu pindahkan ke migration terpisah jika production
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'discount_type') THEN
				CREATE TYPE discount_type AS ENUM ('percent', 'fixed');
			END IF;
		END$$;
	`).Error; err != nil {
		log.Fatalf("❌ Gagal membuat ENUM 'discount_type': %v", err)
	}

	// AutoMigrate hanya saat development (opsional)
	if appEnv != "release" {
		if err := db.AutoMigrate(
			&entity.UserBusiness{},
			&entity.User{},
			&entity.BusinessType{},
			&entity.Business{},
			&entity.Membership{},
			&entity.BusinessBranch{},
			&entity.ProductCategory{},
			&entity.Unit{},
			&entity.Tax{},
			&entity.Discount{},
			&entity.Product{}, // ✅ pastikan ini duluan
			&entity.ProductVariant{},
			&entity.ProductAttribute{},
			&entity.Bundle{},
			&entity.BundleItem{},
			&entity.PaymentMethod{},
			&entity.Customer{},
			&entity.Transaction{},
			&entity.TransactionItem{},
			&entity.TransactionItemAttribute{},
			&entity.Promo{}, // ✅ setelah Product
			&entity.PromoRequiredProduct{},
			&entity.ProductPromo{},
		); err != nil {
			log.Fatalf("❌ AutoMigrate gagal: %v", err)
		}
	} else {
		log.Println("ℹ️ Production mode terdeteksi, AutoMigrate dilewati.")
	}

	return db
}

// CloseDatabaseConnection menutup koneksi database.
func CloseDatabaseConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan koneksi SQL: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Gagal menutup koneksi database: %v", err)
	}
}
