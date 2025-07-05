package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func DeterminePromoType(amount float64) bool {
	return amount <= 1.0
}

func GenerateSKU(name string) string {
	// Ambil 3 huruf pertama dari nama produk
	prefix := strings.ToUpper(name)
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}

	timestamp := time.Now().UnixNano() / 1e6 // millisec
	randomPart := rand.Intn(1000)            // 3 digit acak

	return fmt.Sprintf("%s-%d-%03d", prefix, timestamp, randomPart)
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

// HashOTP meng-hash OTP untuk disimpan dengan aman
func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(hash[:])
}
