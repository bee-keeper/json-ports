package application

import "github.com/bee-keeper/json-ports/internal/domain"

// PortService - intermediary service to handle port related use cases
type PortService struct {
	Repo domain.PortRepository
}

// NewPortService creates a new PortService instance
func NewPortService(repo domain.PortRepository) *PortService {
	return &PortService{
		Repo: repo,
	}
}

// UpsertPort upserts the port information
func (s *PortService) UpsertPort(port *domain.Port) error {
	return s.Repo.UpsertPort(port) // Delegate upsert operation to the repository
}
