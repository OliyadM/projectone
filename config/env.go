package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads variables from the .env file into the system
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system environment variables.")
	}
}

// GetEnv returns an environment variable or a fallback value
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
