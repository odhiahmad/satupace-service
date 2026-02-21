package middleware

import (
	"log"
	"net/http"
	"strings"

	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.BuildErrorResponse(
				"Unauthorized", "INVALID_TOKEN", "user_id", "user_id is not string (UUID)", nil,
			))
			return
		}

		if businessID, ok := claims["business_id"].(string); ok {
			c.Set("business_id", businessID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.BuildErrorResponse(
				"Unauthorized", "INVALID_TOKEN", "business_id", "business_id is not string (UUID)", nil,
			))
			return
		}

		if claims["role_id"] != nil {
			if roleID, ok := claims["role_id"].(float64); ok {
				c.Set("role_id", int(roleID))
			}
		}

		if claims["email"] != nil {
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
		}

		c.Next()
	}
}
