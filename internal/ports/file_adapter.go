package ports

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bee-keeper/json-ports/internal/application"
	"github.com/bee-keeper/json-ports/internal/domain"
)

// FileAdapter handles reading ports from a JSON file
type FileAdapter struct {
	PortService *application.PortService
}

// NewFileAdapter creates a new instance of FileAdapter
func NewFileAdapter(portService *application.PortService) *FileAdapter {
	return &FileAdapter{
		PortService: portService,
	}
}

// UpsertPorts reads ports from a file and upserts them using the service
func (a *FileAdapter) UpsertPorts(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Set up graceful shutdown with OS signal handling
	sigChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to wait for shutdown signals
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %s, shutting down...", sig)
		done <- true
	}()

	// Read the opening brace of the JSON object
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading start of JSON object: %w", err)
	}

	for {
		// Read next token (the key, which is the 'unloc' identifier)
		key, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("error reading port key: %w", err)
		}

		if key == json.Delim('}') {
			break
		}

		// Parse the key ('unloc' code)
		unloc, ok := key.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %T", key)
		}

		// Decode the port data into the Port struct
		var port domain.Port
		if err := decoder.Decode(&port); err != nil {
			return fmt.Errorf("error decoding port data for Unloc %s: %w", unloc, err)
		}

		// Set the Unloc value manually
		port.Unloc = unloc

		// Upsert the port data via the PortService
		if err := a.PortService.UpsertPort(&port); err != nil {
			log.Printf("Error upserting port with Unloc %s: %v", unloc, err)
		} else {
			log.Printf("Successfully upserted port with Unloc %s", unloc)
		}

		// Check if graceful shutdown is requested
		select {
		case <-done:
			log.Println("Graceful shutdown requested, stopping port processing.")
			return nil
		default:
		}

	}
	return nil
}
