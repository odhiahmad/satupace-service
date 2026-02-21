package middleware

import (
	"log"
	"net/http"
	"strings"

	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

		// Ensure this is an access token, not a refresh token
		if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.BuildErrorResponse(
				"Unauthorized", "INVALID_TOKEN_TYPE", "Authorization", "Gunakan access token, bukan refresh token", nil,
			))
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.BuildErrorResponse(
				"Unauthorized", "INVALID_TOKEN", "user_id", "user_id is not string (UUID)", nil,
			))
			return
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.BuildErrorResponse(
				"Unauthorized", "INVALID_TOKEN", "user_id", "user_id is not a valid UUID", nil,
			))
			return
		}

		// Store both string and UUID versions
		c.Set("user_id", userUUID)
		c.Set("user_id_str", userIDStr)

		// Check if user is verified
		if isVerified, ok := claims["is_verified"].(bool); !ok || !isVerified {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.BuildErrorResponse(
				"Forbidden", "NOT_VERIFIED", "user", "Akun belum diverifikasi", nil,
			))
			return
		}

		// Check if user is active
		if isActive, ok := claims["is_active"].(bool); !ok || !isActive {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.BuildErrorResponse(
				"Forbidden", "ACCOUNT_INACTIVE", "user", "Akun tidak aktif", nil,
			))
			return
		}

		if claims["email"] != nil {
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
		}

		if claims["phone_number"] != nil {
			if phoneNumber, ok := claims["phone_number"].(string); ok {
				c.Set("phone_number", phoneNumber)
			}
		}

		c.Next()
	}
}
