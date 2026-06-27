package router

import (
	"github.com/gin-gonic/gin"
	"github.com/areyoush/surfspace/internal/auth"
	"github.com/areyoush/surfspace/internal/links"
	"github.com/areyoush/surfspace/internal/middleware"
)

func Setup(r *gin.Engine, authHandler *auth.Handler, linksHandler *links.Handler, authRepo *auth.Repository, jwtSecret string) {
	v1 := r.Group("/api/v1")
	auth.RegisterRoutes(v1.Group("/auth"), authHandler)

	protected := v1.Group("/")
	protected.Use(middleware.Auth(jwtSecret, authRepo))

	links.RegisterRoutes(protected, linksHandler, r.Group("/"))
}