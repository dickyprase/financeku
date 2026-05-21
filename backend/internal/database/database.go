package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/financeku/backend/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations executed successfully")
	return nil
}

func RunSeed(db *sql.DB) error {
	_, err := db.Exec(seedSQL)
	if err != nil {
		return fmt.Errorf("failed to run seed: %w", err)
	}
	log.Println("Seed data inserted successfully")
	return nil
}
