package router

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/areyoush/surfspace/internal/auth"
	"github.com/areyoush/surfspace/internal/links"
	"github.com/areyoush/surfspace/internal/middleware"
	"github.com/areyoush/surfspace/internal/ratelimit"
)

func Setup(r *gin.Engine, authHandler *auth.Handler, linksHandler *links.Handler, authRepo *auth.Repository, redisClient *redis.Client, jwtSecret string) {
	rl := ratelimit.NewRateLimiter(redisClient, 60, time.Minute)

	v1 := r.Group("/api/v1")
	v1.Group("/auth").Use(rl.RateLimitByIP())
	auth.RegisterRoutes(v1.Group("/auth"), authHandler)

	protected := v1.Group("/")
	protected.Use(middleware.Auth(jwtSecret, authRepo))
	protected.Use(rl.RateLimitByUserID())
	publicRg := r.Group("/")
	publicRg.Use(rl.RateLimitByIP())
	links.RegisterRoutes(protected, linksHandler, publicRg)
}