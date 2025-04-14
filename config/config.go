package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	LogLevel   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Invalid or no .env file")
	}

	appEnv := getEnv("APP_ENV", "local")

	if appEnv == "docker" {
		log.Println("Running in Docker mode, ignoring .env file")
		return &Config{
			DBHost:     getEnv("DB_HOST", "db"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", "postgres"),
			DBName:     getEnv("DB_NAME", "pvz"),
			JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
			LogLevel:   getEnv("LOG_LEVEL", "debug"),
		}
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "pvz"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
		LogLevel:   getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
