package usecases

import (
	"fmt"

	"kodkafa/internal/app/dto"
	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
)

// DeletePluginUseCase handles removing a plugin and its associated data.
type DeletePluginUseCase struct {
	pluginRepo ports.PluginRepository
	stateStore ports.StateStore
	usageStore ports.UsageStore
	installer  ports.DependencyInstaller
}

// NewDeletePluginUseCase creates a new DeletePluginUseCase.
func NewDeletePluginUseCase(
	pluginRepo ports.PluginRepository,
	stateStore ports.StateStore,
	usageStore ports.UsageStore,
	installer ports.DependencyInstaller,
) *DeletePluginUseCase {
	return &DeletePluginUseCase{
		pluginRepo: pluginRepo,
		stateStore: stateStore,
		usageStore: usageStore,
		installer:  installer,
	}
}

// DeletePluginInput represents the input for DeletePluginUseCase.
type DeletePluginInput struct {
	PluginName string
	RemoveDeps bool
}

// Execute removes a plugin and cleans up its state and usage stats.
func (uc *DeletePluginUseCase) Execute(input DeletePluginInput) (dto.DeletePluginResult, error) {
	if input.PluginName == "" {
		return dto.DeletePluginResult{Success: false, Message: "plugin name is required"}, fmt.Errorf("plugin name is required")
	}

	result := dto.DeletePluginResult{
		PluginName: input.PluginName,
	}

	// 1. Check if plugin exists
	exists, err := uc.pluginRepo.Exists(input.PluginName)
	if err != nil {
		result.Message = err.Error()
		return result, err
	}
	if !exists {
		result.Message = "plugin not found"
		return result, fmt.Errorf("plugin not found: %s", input.PluginName)
	}

	// 2. Remove dependencies if requested (must be done before Remove)
	if input.RemoveDeps {
		// Get plugin entity for its metadata
		plugin, _ := uc.pluginRepo.Get(input.PluginName)
		if plugin != nil {
			if err := uc.installer.Uninstall(plugin); err != nil {
				result.Message = fmt.Sprintf("shared dependency cleanup failed: %v", err)
			}
		}

		if err := uc.pluginRepo.RemoveDeps(input.PluginName); err != nil {
			// Log warning or record in result message
			result.Message = fmt.Sprintf("local dependency cleanup failed: %v", err)
		}
	}

	// 2.1 Delete plugin from repository (metadata and binary)
	if err := uc.pluginRepo.Remove(input.PluginName); err != nil {
		result.Message = fmt.Errorf("failed to remove plugin: %w", err).Error()
		return result, err
	}

	// 3. Delete plugin state
	if err := uc.stateStore.Delete(input.PluginName); err != nil {
		// Log warning but continue
	}

	// 4. Update usage stats
	usage, err := uc.usageStore.Read()
	if err == nil && usage != nil {
		updated := false

		// Remove from RecentlyUsed
		newRecent := make([]entities.UsageEntry, 0)
		for _, entry := range usage.RecentlyUsed {
			if entry.PluginName != input.PluginName {
				newRecent = append(newRecent, entry)
			} else {
				updated = true
			}
		}
		usage.RecentlyUsed = newRecent

		// Remove from MostUsed
		newMost := make([]entities.UsageEntry, 0)
		for _, entry := range usage.MostUsed {
			if entry.PluginName != input.PluginName {
				newMost = append(newMost, entry)
			} else {
				updated = true
			}
		}
		usage.MostUsed = newMost

		if updated {
			_ = uc.usageStore.Write(usage)
		}
	}

	result.Success = true
	result.Message = "plugin deleted successfully"
	return result, nil
}
