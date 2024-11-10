package infra

import (
	"fmt"

	"github.com/bee-keeper/json-ports/internal/domain"
	"gorm.io/gorm"
)

// MigrateDB handles database schema migrations
func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Port{}); err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}
	return nil
}
