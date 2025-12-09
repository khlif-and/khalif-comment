package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

)

type RateLimitConfig struct {
	Limit  int
	Window time.Duration
}

func RateLimit(rdb *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rl:%s", ip)

		count, err := rdb.Incr(c.Request.Context(), key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(c.Request.Context(), key, config.Window)
		}

		remaining := int64(config.Limit) - count
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))

		if count > int64(config.Limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}

		c.Next()
	}
}