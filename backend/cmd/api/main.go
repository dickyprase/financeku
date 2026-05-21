package main

import (
	"log"
	"net/http"
	"os"

	"github.com/financeku/backend/internal/config"
	"github.com/financeku/backend/internal/database"
	"github.com/financeku/backend/internal/router"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Run seed (only if --seed flag is passed)
	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		if err := database.RunSeed(db); err != nil {
			log.Fatalf("Failed to run seed: %v", err)
		}
	}

	// Setup router
	handler := router.Setup(db, cfg)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("FinanceKu API server starting on %s", addr)
	log.Printf("Environment: %s", getEnv("APP_ENV", "development"))

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
