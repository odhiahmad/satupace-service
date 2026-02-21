package middleware

import (
	"net/http"

	"run-sync/helper"
	"run-sync/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProfileRequired blocks access to features that require a runner profile setup.
// Must be used AFTER AuthorizeJWT middleware.
func ProfileRequired(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.MustGet("user_id").(uuid.UUID)

		user, err := userRepo.FindById(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.BuildErrorResponse(
				"Unauthorized", "USER_NOT_FOUND", "user", "User tidak ditemukan", nil,
			))
			return
		}

		if !user.HasProfile {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.BuildErrorResponse(
				"Profil belum lengkap",
				"PROFILE_REQUIRED",
				"user",
				"Silakan lengkapi profil runner Anda terlebih dahulu sebelum mengakses fitur ini",
				nil,
			))
			return
		}

		c.Next()
	}
}
