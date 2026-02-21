package service

import (
	"fmt"
	"os"
	"time"

	"run-sync/helper"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userId string, phoneNumber string, email *string, expiredAt time.Time) string
	GenerateRefreshToken(userId string) string
	ValidateToken(token string) (*jwt.Token, error)
	ValidateRefreshToken(token string) (*jwt.Token, error)
}

type jwtCustomClaims struct {
	UserId      string  `json:"user_id"`
	PhoneNumber string  `json:"phone_number"`
	Email       *string `json:"email"`
	IsVerified  bool    `json:"is_verified"`
	IsActive    bool    `json:"is_active"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey        string
	refreshSecretKey string
	issuer           string
}

func NewJwtService() JWTService {
	return &jwtService{
		issuer:           "run-sync",
		secretKey:        getSecretKey(),
		refreshSecretKey: getRefreshSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "sdfnkjsdf28fmn*(&^%^%&bjsdfgsQ$@$sadjfsdfx"
	}
	return secretKey
}

func getRefreshSecretKey() string {
	secretKey := os.Getenv("JWT_REFRESH_SECRET")
	if secretKey == "" {
		secretKey = "rf_sdfnkjsdf28fmn*(&^%^%&bjsdfgsQ$@$sadjfsdfx"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(userId string, phoneNumber string, email *string, expiredAt time.Time) string {
	claims := jwt.MapClaims{
		"user_id":      userId,
		"phone_number": phoneNumber,
		"email":        email,
		"is_verified":  true, // Only verified users should get tokens
		"is_active":    true, // Only active users should get tokens
		"token_type":   "access",
		"exp":          expiredAt.Unix(),
		"iss":          j.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	helper.ErrorPanic(err)

	return signedToken
}

func (j *jwtService) GenerateRefreshToken(userId string) string {
	claims := jwt.MapClaims{
		"user_id":    userId,
		"token_type": "refresh",
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"iss":        j.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.refreshSecretKey))
	helper.ErrorPanic(err)

	return signedToken
}

func (j *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	return token, err
}

func (j *jwtService) ValidateRefreshToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.refreshSecretKey), nil
	})

	return token, err
}
