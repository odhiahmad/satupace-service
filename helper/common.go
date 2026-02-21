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
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

// GenerateOTPCode generates a random OTP code of specified length
func GenerateOTPCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

// HashOTP creates a SHA256 hash of the OTP for secure storage
func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(hash[:])
}

// ExtractPublicIDFromURL extracts Cloudinary public ID from a URL
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
	publicID = strings.TrimSuffix(publicID, filepath.Ext(publicID))
	return publicID, nil
}

// DeleteFromCloudinary deletes an asset from Cloudinary by public ID
func DeleteFromCloudinary(publicID string) error {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return err
	}

	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// LowerStringPtr converts a string pointer to lowercase
func LowerStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	lowered := strings.ToLower(*s)
	return &lowered
}

// StringValue safely extracts value from string pointer
func StringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// StringPtr creates a pointer to a string value
func StringPtr(s string) *string {
	return &s
}

// HashAndSalt securely hashes a password with bcrypt
func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}

// ComparePassword verifies a password against its hash
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPassword hashes a password string
func HashPassword(password string) string {
	return HashAndSalt([]byte(password))
}

// IsEmail validates if a string is a valid email format
func IsEmail(input string) bool {
	return emailRegex.MatchString(input)
}
