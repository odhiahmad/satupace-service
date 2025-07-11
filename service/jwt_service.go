package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
)

type JWTService interface {
	GenerateToken(user entity.UserBusiness) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJwtService() JWTService {
	return &jwtService{
		issuer:    "loka",
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

func (j *jwtService) GenerateToken(user entity.UserBusiness) string {
	claims := jwt.MapClaims{
		"user_id":      user.Id,
		"phone_number": user.PhoneNumber,
		"business_id":  user.BusinessId,
		"email":        user.Email,
		"role_id":      user.RoleId,
		"exp":          time.Now().Add(100 * 365 * 24 * time.Hour).Unix(), // tidak expired dalam waktu dekat
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
