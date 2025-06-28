package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response := helper.BuildErrorResponse(
				"Unauthorized",
				"UNAUTHORIZED",
				"Authorization",
				"No valid Bearer token found",
				nil,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			log.Println("JWT error:", err)
			response := helper.BuildErrorResponse(
				"Unauthorized",
				"INVALID_TOKEN",
				"Authorization",
				err.Error(),
				nil,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Simpan klaim ke context jika diperlukan
		claims := token.Claims.(jwt.MapClaims)
		log.Println("Claim[userId]:", claims["userId"])
		log.Println("Claim[issuer]:", claims["issuer"])

		// Simpan userId ke context jika ingin digunakan di handler
		c.Set("userId", claims["userId"])
		c.Next()
	}
}
