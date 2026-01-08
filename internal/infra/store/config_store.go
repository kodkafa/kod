package store

import (
	"kodkafa/internal/domain/ports"
)

// ConfigStoreImpl implements ports.ConfigStore using JSON files.
type ConfigStoreImpl struct {
	store *JSONStore
}

// NewConfigStore creates a new ConfigStore implementation.
func NewConfigStore(baseDir string) ports.ConfigStore {
	return &ConfigStoreImpl{
		store: NewJSONStore(baseDir),
	}
}

// Read reads the global configuration.
func (cs *ConfigStoreImpl) Read() (*ports.Config, error) {
	var config ports.Config
	err := cs.store.Read("config.json", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Write persists the global configuration.
func (cs *ConfigStoreImpl) Write(config *ports.Config) error {
	return cs.store.Write("config.json", config)
}
