package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func DeterminePromoType(amount float64) string {
	if amount <= 1.0 {
		return "percent"
	}
	return "fixed"
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
