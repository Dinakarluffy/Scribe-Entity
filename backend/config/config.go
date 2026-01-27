package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// ProjectRoot holds the absolute path to the project root
// Example: D:\INCUBRIX PROJECT\Scribe-Entity-Classification
var ProjectRoot string

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() *Config {
	// Start from current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	var envPath string

	// Walk upwards until .env is found
	for {
		tryPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(tryPath); err == nil {
			envPath = tryPath
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}

	if envPath == "" {
		log.Println("No .env file found, relying on system environment variables")
	} else {
		if err := godotenv.Load(envPath); err != nil {
			log.Println("Failed to load .env at", envPath)
		} else {
			log.Println(".env loaded from", envPath)

			// 🔑 VERY IMPORTANT
			// This is used to resolve Python-Worker path correctly
			ProjectRoot = filepath.Dir(envPath)
		}
	}
		return &Config{
			Port: getEnv("PORT", "8080"),

			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", ""),
			DBName:     getEnv("DB_NAME", "entity_classification"),
			DBSSLMode:  getEnv("DB_SSLMODE", "require"),
		}
	}

	func getEnv(key, fallback string) string {
		if value, exists := os.LookupEnv(key); exists && value != "" {
			return value
		}
		return fallback
	}
