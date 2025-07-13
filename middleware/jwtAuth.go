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

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("business_id", int(claims["business_id"].(float64)))
		if claims["branch_id"] != nil {
			c.Set("branch_id", int(claims["branch_id"].(float64)))
		}
		if claims["role_id"] != nil {
			c.Set("role_id", int(claims["role_id"].(float64)))
		}
		if claims["email"] != nil {
			c.Set("email", claims["email"].(string))
		}

		c.Next()
	}
}
