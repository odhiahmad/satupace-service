package helper

import (
	"fmt"
	"time"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

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
