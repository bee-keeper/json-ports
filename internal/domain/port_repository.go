package domain

// PortRepository defines the interface for repository operations
type PortRepository interface {
	UpsertPort(port *Port) error
}
