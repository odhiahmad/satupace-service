package helper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleNotFound(ctx *gin.Context, err error, message string) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": message,
		})
		return true
	}
	return false
}
