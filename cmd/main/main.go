package main

import (
	"fmt"
	"log"

	"github.com/dorik33/Test/internal/api"
	"github.com/dorik33/Test/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Music API
// @version 1.0
// @description API for managing songs

// @host localhost:8080
// @BasePath /

func applyMigrations(cfg *config.Config) error {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		),
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := applyMigrations(cfg); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	apiServer := api.New(cfg)
	if err := apiServer.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
