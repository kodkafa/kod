package ports

import "kodkafa/internal/domain/entities"

// UsageStore defines the interface for global usage statistics persistence.
type UsageStore interface {
	// Read reads the global usage statistics.
	Read() (*entities.UsageStats, error)
	// Write persists the global usage statistics.
	Write(stats *entities.UsageStats) error
}
