package main

import (
	"log"

	"github.com/bee-keeper/json-ports/internal/application"
	"github.com/bee-keeper/json-ports/internal/infra"
	"github.com/bee-keeper/json-ports/internal/ports"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialise SQLite in-memory DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := infra.MigrateDB(db); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Initialise SQLite repository adapter
	portRepo := infra.NewPortRepositorySQLite(db)

	// Initialise application service
	portService := application.NewPortService(portRepo)

	// Initialise file adapter to read ports from JSON file
	fileAdapter := ports.NewFileAdapter(portService)

	// Read and upsert ports from the file
	if err := fileAdapter.UpsertPorts("./data/ports.json"); err != nil {
		log.Fatalf("Error reading and upserting ports: %v", err)
	}
}
