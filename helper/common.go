package helper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

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
