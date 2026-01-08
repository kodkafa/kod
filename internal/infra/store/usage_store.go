package store

import (
	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
)

// UsageStoreImpl implements ports.UsageStore using JSON files.
type UsageStoreImpl struct {
	store *JSONStore
}

// NewUsageStore creates a new UsageStore implementation.
func NewUsageStore(baseDir string) ports.UsageStore {
	return &UsageStoreImpl{
		store: NewJSONStore(baseDir),
	}
}

// Read reads the global usage statistics.
func (us *UsageStoreImpl) Read() (*entities.UsageStats, error) {
	var stats entities.UsageStats
	err := us.store.Read("usage.json", &stats)
	if err != nil {
		// If file doesn't exist, create default stats
		if !us.store.Exists("usage.json") {
			newStats := entities.NewUsageStats()
			if writeErr := us.Write(newStats); writeErr != nil {
				return nil, writeErr
			}
			return newStats, nil
		}
		return nil, err
	}
	return &stats, nil
}

// Write persists the global usage statistics.
func (us *UsageStoreImpl) Write(stats *entities.UsageStats) error {
	return us.store.Write("usage.json", stats)
}
