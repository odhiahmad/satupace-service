package config

import (
	"fmt"
	"log"
	"os"

	"run-sync/entity"

	"github.com/joho/godotenv"
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
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}

	if ginMode != "release" {
		if err := db.AutoMigrate(
			&entity.User{},
			&entity.RunnerProfile{},
			&entity.RunGroup{},
			&entity.RunGroupMember{},
			&entity.RunActivity{},
			&entity.DirectMatch{},
			&entity.DirectChatMessage{},
			&entity.GroupChatMessage{},
			&entity.UserPhoto{},
			&entity.SafetyLog{},
			&entity.UserBiometric{},
			&entity.Notification{},
			&entity.UserDeviceToken{},
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
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke DB WhatsApp: %v", err)
	}

	return db
}
