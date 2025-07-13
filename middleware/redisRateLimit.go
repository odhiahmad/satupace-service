package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/helper"
)

func RateLimit(redisHelper *helper.RedisHelper, max int, period time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		err := redisHelper.AllowRequest("rate_limit:"+ip, max, period)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, helper.BuildErrorResponse(
				"Terlalu banyak permintaan, silakan coba lagi nanti.",
				"TOO_MANY_REQUESTS",
				"throttle",
				"Terlalu banyak permintaan ke server dalam waktu singkat.",
				helper.EmptyObj{},
			))
			return
		}
		c.Next()
	}
}
