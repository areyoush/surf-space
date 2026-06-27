package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/areyoush/surfspace/internal/auth"
	"github.com/areyoush/surfspace/internal/config"
	"github.com/areyoush/surfspace/internal/db"
	"github.com/areyoush/surfspace/internal/router"
	"github.com/areyoush/surfspace/internal/links"
	"github.com/areyoush/surfspace/internal/cache"
	
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	database := db.Connect(cfg.DBDSN)
	db.RunMigrations(database)
	authRepo := auth.NewRepository(database)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)
	linksRepo := links.NewRepository(database)
	linksSvc := links.NewService(linksRepo)
	linksHandler := links.NewHandler(linksSvc)
	redisClient := cache.New(cfg.RedisAddr, cfg.RedisPassword)

	r := gin.Default()
	router.Setup(r, authHandler, linksHandler, authRepo, redisClient, cfg.JWTSecret)
	log.Println("Server starting on port", cfg.Port)

	r.Run(":" + cfg.Port)
	

	
}