package ports

import "kodkafa/internal/domain/entities"

// StateStore defines the interface for per-plugin state persistence.
type StateStore interface {
	// Read reads the state for a plugin, creating it if it doesn't exist.
	Read(pluginName string) (*entities.PluginState, error)
	// Write persists the state for a plugin.
	Write(state *entities.PluginState) error
	// Delete removes the state for a plugin.
	Delete(pluginName string) error
}
