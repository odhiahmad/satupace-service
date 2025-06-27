package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/odhiahmad/kasirku-service/helper"
)

type JWTService interface {
	GenerateToken(userId int) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserdId string `json:"userId"`
	jwt.RegisteredClaims
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

func (j *jwtService) GenerateToken(UserId int) string {
	claims := jwt.MapClaims{}
	claims["userId"] = UserId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["issuer"] = "loka"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	helper.ErrorPanic(err)

	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
