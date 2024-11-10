package infra

import (
	"github.com/bee-keeper/json-ports/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PortRepositorySQLite is the SQLite implementation of the PortRepository
type PortRepositorySQLite struct {
	DB *gorm.DB
}

// NewPortRepositorySQLite creates a new PortRepositorySQLite instance
func NewPortRepositorySQLite(db *gorm.DB) *PortRepositorySQLite {
	return &PortRepositorySQLite{
		DB: db,
	}
}

// UpsertPort receiver binds UpsertPort to repo
func (r *PortRepositorySQLite) UpsertPort(port *domain.Port) error {
	if err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unloc"}},
		UpdateAll: true,
	}).Create(port).Error; err != nil {
		return err
	}
	return nil
}
