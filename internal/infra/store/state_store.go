package store

import (
	"fmt"
	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
	"os"
)

// StateStoreImpl implements ports.StateStore using JSON files.
type StateStoreImpl struct {
	store   *JSONStore
	baseDir string
}

// NewStateStore creates a new StateStore implementation.
func NewStateStore(baseDir string) ports.StateStore {
	return &StateStoreImpl{
		store:   NewJSONStore(baseDir),
		baseDir: baseDir,
	}
}

// Read reads the state for a plugin, creating it if it doesn't exist.
func (ss *StateStoreImpl) Read(pluginName string) (*entities.PluginState, error) {
	filename := fmt.Sprintf("state/%s.json", pluginName)

	var state entities.PluginState
	err := ss.store.Read(filename, &state)
	if err != nil {
		// If file doesn't exist, create a new state
		if !ss.store.Exists(filename) {
			newState := entities.NewPluginState(pluginName)
			if writeErr := ss.Write(newState); writeErr != nil {
				return nil, fmt.Errorf("failed to create initial state: %w", writeErr)
			}
			return newState, nil
		}
		return nil, fmt.Errorf("failed to read state: %w", err)
	}

	// Ensure plugin name matches
	state.PluginName = pluginName
	return &state, nil
}

// Write persists the state for a plugin.
func (ss *StateStoreImpl) Write(state *entities.PluginState) error {
	filename := fmt.Sprintf("state/%s.json", state.PluginName)
	return ss.store.Write(filename, state)
}

// Delete removes the state for a plugin.
func (ss *StateStoreImpl) Delete(pluginName string) error {
	path := fmt.Sprintf("%s/state/%s.json", ss.baseDir, pluginName)
	return os.Remove(path)
}
