package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{client: client, limit: limit, window: window}
}

func (rl *RateLimiter) RateLimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:ip:%s", c.ClientIP())
		if !rl.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded, try again later"})
			return
		}
		c.Next()
	}
}

func (rl *RateLimiter) RateLimitByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		key := fmt.Sprintf("ratelimit:user:%s", userID.(string))
		if !rl.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded, try again later"})
			return
		}
		c.Next()
	}
}

func (rl *RateLimiter) allow(key string) bool {
	ctx := context.Background()
	now := time.Now().UnixNano()
	windowStart := now - rl.window.Nanoseconds()

	pipe := rl.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
	pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, rl.window)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return true
	}

	count := results[2].(*redis.IntCmd).Val()
	return count <= int64(rl.limit)
}