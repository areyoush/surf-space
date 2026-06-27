package config

import (
	"fmt"
	"os"
	"log"
)


type Config struct {
	Port		string
	DBHost		string
	DBPort		string
	DBUser		string
	DBPassword	string
	DBName		string
	DBSSLMode	string
	DBDSN		string
	JWTSecret	string
}


func Load() *Config {
	cfg := &Config{
		Port:		os.Getenv("PORT"),
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        DBSSLMode:  os.Getenv("DB_SSLMODE"),
        JWTSecret:  os.Getenv("JWT_SECRET"),
	}
    

	if cfg.Port == "" {
    	cfg.Port = "8080"
	}

	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "disable"
	}
	
	if cfg.DBHost == "" {
		log.Fatal("DB_HOST environment variable is not set")
	}

	if cfg.DBPort == "" {
    	cfg.DBPort = "5432"
	}
	
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	cfg.DBDSN = fmt.Sprintf(
    	"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
     	cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	return cfg
	
}