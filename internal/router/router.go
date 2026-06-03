package router

import (

	"github.com/gin-gonic/gin"
	"github.com/areyoush/surfspace/internal/auth"
)

func Setup(r *gin.Engine, authHandler *auth.Handler) {
	v1 := r.Group("/api/v1")
	auth.RegisterRoutes(v1.Group("/auth"), authHandler)
}

