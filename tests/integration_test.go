package tests

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/bee-keeper/json-ports/internal/application"
	"github.com/bee-keeper/json-ports/internal/domain"
	"github.com/bee-keeper/json-ports/internal/infra"
	"github.com/bee-keeper/json-ports/internal/ports"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initialises an in-memory SQLite database for testing
func setupTestDB(t *testing.T) (*gorm.DB, *infra.PortRepositorySQLite) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Port{}); err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	repo := infra.NewPortRepositorySQLite(db)
	return db, repo
}

// TestImportPortFromFile reads the single_port.json file, imports the port, and verifies the data in the DB
func TestImportPortFromFile(t *testing.T) {
	db, repo := setupTestDB(t)
	portService := application.NewPortService(repo)

	fileAdapter := ports.NewFileAdapter(portService)
	filePath := "./data/single_port.json"

	err := fileAdapter.UpsertPorts(filePath)
	assert.NoError(t, err)

	expectedPort := &domain.Port{
		Unloc:    "AEAJM",
		Name:     "Ajman",
		City:     "Ajman",
		Country:  "United Arab Emirates",
		Province: "Ajman",
		Timezone: "Asia/Dubai",
		Coordinates: datatypes.JSON(
			[]byte(`[55.5136433,25.4052165]`),
		),
		Code:   "52000",
		Unlocs: datatypes.JSON([]byte(`["AEAJM"]`)),
	}

	// Query the database to verify the port was inserted and check the fields
	var insertedPort domain.Port
	err = db.First(&insertedPort, "unloc = ?", expectedPort.Unloc).Error
	assert.NoError(t, err)

	// Assert that the fields in the database match the expected values
	assert.Equal(t, expectedPort.Unloc, insertedPort.Unloc)
	assert.Equal(t, expectedPort.Name, insertedPort.Name)
	assert.Equal(t, expectedPort.City, insertedPort.City)
	assert.Equal(t, expectedPort.Country, insertedPort.Country)
	assert.Equal(t, expectedPort.Province, insertedPort.Province)
	assert.Equal(t, expectedPort.Timezone, insertedPort.Timezone)
	assert.Equal(t, expectedPort.Code, insertedPort.Code)

	insertedUnlocs, _ := json.Marshal(insertedPort.Unlocs)
	assert.Equal(t, string(expectedPort.Unlocs), string(insertedUnlocs))

	insertedCoordinates, _ := json.Marshal(insertedPort.Coordinates)
	assert.Equal(t, string(expectedPort.Coordinates), string(insertedCoordinates))

}

// TestUpdateCoordinatesAfterImport tests that importing a port from a file updates an existing port
func TestUpdateCoordinatesAfterImport(t *testing.T) {
	db, repo := setupTestDB(t)
	portService := application.NewPortService(repo)

	initialPort := &domain.Port{
		Unloc:    "AEAJM",
		Name:     "Ajman",
		City:     "Ajman",
		Country:  "United Arab Emirates",
		Province: "Ajman",
		Timezone: "Asia/Dubai",
		Coordinates: datatypes.JSON(
			[]byte(`[56.5136433,26.4052165]`),
		),
		Code:   "52000",
		Unlocs: datatypes.JSON([]byte(`["AEAJM"]`)),
	}

	// Insert the initial port into the database
	err := repo.UpsertPort(initialPort)
	assert.NoError(t, err)

	// Create a FileAdapter to read and import the ports from the JSON file
	fileAdapter := ports.NewFileAdapter(portService)
	filePath := "./data/single_port.json"

	// Import the port from the file, which contains updated data
	err = fileAdapter.UpsertPorts(filePath)
	assert.NoError(t, err)

	// Query the database to verify the port was updated
	var updatedPort domain.Port
	err = db.First(&updatedPort, "unloc = ?", initialPort.Unloc).Error
	assert.NoError(t, err)

	// Assert that the coordinates were updated to match the data in the file
	expectedCoordinates := datatypes.JSON([]byte(`[55.5136433,25.4052165]`))
	assert.JSONEq(t, string(expectedCoordinates), string(updatedPort.Coordinates))

	// Assert other fields remain unchanged
	assert.Equal(t, initialPort.Name, updatedPort.Name)
	assert.Equal(t, initialPort.City, updatedPort.City)
	assert.Equal(t, initialPort.Country, updatedPort.Country)
	assert.Equal(t, initialPort.Province, updatedPort.Province)
	assert.Equal(t, initialPort.Timezone, updatedPort.Timezone)
	assert.Equal(t, initialPort.Code, updatedPort.Code)

	insertedUnlocs, _ := json.Marshal(updatedPort.Unlocs)
	assert.JSONEq(t, string(initialPort.Unlocs), string(insertedUnlocs))
}

// TestImportInvalidPort tests importing invalid port data from a file
func TestImportInvalidPort(t *testing.T) {
	_, repo := setupTestDB(t)
	portService := application.NewPortService(repo)

	fileAdapter := ports.NewFileAdapter(portService)
	filePath := "./data/invalid_port.json"

	// Attempt to import from the invalid file
	err := fileAdapter.UpsertPorts(filePath)
	log.Println(err)
	assert.Error(t, err, "Expected an error when importing invalid port data")
	assert.Contains(t, err.Error(), "error decoding port data for Unloc AEAJM: invalid character ']' looking for beginning of object key string")
}
