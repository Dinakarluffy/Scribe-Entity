package main

import (
	"log"
	"net/http"
	"os"

	"your_project/config"
	"your_project/repository"
	"your_project/routes"
)

func main() {
	// Load config (this also loads .env internally)
	cfg := config.Load()

	// Debug (now this WILL work)
	log.Println("PYTHON_WORKER_PATH =", os.Getenv("PYTHON_WORKER_PATH"))

	// Init DB
	db := config.InitDB(cfg)
	repository.InitRepository(db)

	// Init routes
	router := routes.RegisterRoutes()

	// Start server
	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
