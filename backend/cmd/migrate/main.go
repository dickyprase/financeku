package main

import (
	"fmt"
	"log"
	"os"

	"github.com/financeku/backend/internal/config"
	"github.com/financeku/backend/internal/database"
)

// Standalone migration tool
// Usage: go run migrations/migrate.go [--seed]
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Running migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed successfully!")

	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		fmt.Println("Running seed...")
		if err := database.RunSeed(db); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		fmt.Println("Seed completed successfully!")
	}
}
