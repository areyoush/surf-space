package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/areyoush/surfspace/internal/config"
	"github.com/areyoush/surfspace/internal/db"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	database := db.Connect(cfg.DBDSN)

	_ = database
	r := gin.Default()
	r.Run(":" + cfg.Port)

	log.Println("Server starting on port", cfg.Port)
}