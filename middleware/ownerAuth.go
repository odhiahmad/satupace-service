package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// misalnya role_id = 1 untuk owner
		if roleInt, ok := roleID.(int); !ok || roleInt != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden: owner role required",
			})
			return
		}

		c.Next()
	}
}
