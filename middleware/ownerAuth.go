package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		userClaims, ok := claims.(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		role, ok := userClaims["role"].(string)
		if !ok || role != "owner" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden: owner role required",
			})
			return
		}

		c.Next()
	}
}
