package ports

import "kodkafa/internal/domain/entities"

// PluginRepository defines the interface for plugin storage and retrieval.
type PluginRepository interface {
	// List returns all installed plugins.
	List() ([]entities.Plugin, error)
	// Get returns a plugin by name, or an error if not found.
	Get(name string) (*entities.Plugin, error)
	// Add registers a new plugin from a local path or remote URL.
	Add(source string) (*entities.Plugin, error)
	// Remove removes a plugin by name.
	Remove(name string) error
	// RemoveDeps deletes plugin dependencies
	RemoveDeps(name string) error
	// Exists checks if a plugin with the given name exists.
	Exists(name string) (bool, error)
}
