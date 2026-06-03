package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/areyoush/surfspace/internal/auth"
	"github.com/areyoush/surfspace/internal/config"
	"github.com/areyoush/surfspace/internal/db"
	"github.com/areyoush/surfspace/internal/router"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	database := db.Connect(cfg.DBDSN)
	db.RunMigrations(database)
	authRepo := auth.NewRepository(database)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)

	r := gin.Default()
	router.Setup(r, authHandler)
	log.Println("Server starting on port", cfg.Port)

	r.Run(":" + cfg.Port)
	

	
}