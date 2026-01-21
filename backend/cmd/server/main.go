package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"your_project/config"
	"your_project/repository"
	"your_project/routes"
)

func main() {
	// 🔴 MUST BE FIRST: load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// 🔍 Debug log (keep for now)
	log.Println("PYTHON_WORKER_PATH =", os.Getenv("PYTHON_WORKER_PATH"))

	// Load config AFTER env is loaded
	cfg := config.Load()

	// Init DB
	db := config.InitDB(cfg)
	repository.InitRepository(db)

	// Init routes
	router := routes.RegisterRoutes()

	// Start server
	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
