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
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJwtService() JWTService {
	return &jwtService{
		issuer:    "run-sync",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "sdfnkjsdf28fmn*(&^%^%&bjsdfgsQ$@$sadjfsdfx"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(userId string, phoneNumber string, email *string, expiredAt time.Time) string {
	claims := jwt.MapClaims{
		"user_id":      userId,
		"phone_number": phoneNumber,
		"email":        email,
		"exp":          expiredAt.Unix(),
		"iss":          j.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
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
