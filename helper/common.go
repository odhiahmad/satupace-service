package helper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/odhiahmad/kasirku-service/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func DeterminePromoType(amount float64) bool {
	return amount <= 1.0
}

func GenerateSKU(name string) string {
	prefix := strings.ToUpper(name)
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}

	timestamp := time.Now().Format("0601021504")
	randomPart := rand.Intn(100)

	return fmt.Sprintf("%s%s%02d", prefix, timestamp, randomPart)
}

func GenerateRandomToken(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateOTPCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(hash[:])
}

func ExtractPublicIDFromURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	parts := strings.Split(parsed.Path, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("url path tidak valid")
	}

	index := slices.Index(parts, "upload")
	if index == -1 || index+1 >= len(parts) {
		return "", fmt.Errorf("url tidak mengandung /upload/")
	}

	publicID := strings.Join(parts[index+1:], "/")
	publicID = strings.TrimSuffix(publicID, filepath.Ext(publicID)) // hapus ekstensi
	return publicID, nil
}

func DeleteFromCloudinary(publicID string) error {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	if err != nil {
		return err
	}

	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

func LowerStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	lowered := strings.ToLower(*s)
	return &lowered
}

func StringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func GenerateBillNumber(db *gorm.DB) (string, error) {
	today := time.Now().Format("20060102") // 20250625
	prefix := "TRX-" + today + "-"

	var count int64
	err := db.Model(&entity.Transaction{}).
		Where("DATE(created_at) = ?", time.Now().Format("2006-01-02")).
		Count(&count).Error
	if err != nil {
		return "", err
	}

	billNumber := fmt.Sprintf("%s%04d", prefix, count+1)
	return billNumber, nil
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func IsEmail(input string) bool {
	return emailRegex.MatchString(input)
}
