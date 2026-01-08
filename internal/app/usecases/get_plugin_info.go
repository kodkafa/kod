package usecases

import (
	"fmt"

	"kodkafa/internal/app/dto"
	"kodkafa/internal/domain/ports"
)

// GetPluginInfoUseCase handles fetching detailed information about a plugin.
type GetPluginInfoUseCase struct {
	pluginRepo ports.PluginRepository
	stateStore ports.StateStore
}

// NewGetPluginInfoUseCase creates a new GetPluginInfoUseCase.
func NewGetPluginInfoUseCase(
	pluginRepo ports.PluginRepository,
	stateStore ports.StateStore,
) *GetPluginInfoUseCase {
	return &GetPluginInfoUseCase{
		pluginRepo: pluginRepo,
		stateStore: stateStore,
	}
}

// GetPluginInfoInput represents the input for GetPluginInfoUseCase.
type GetPluginInfoInput struct {
	PluginName   string
	HistoryLimit int
}

// Execute returns detailed information about a plugin.
func (uc *GetPluginInfoUseCase) Execute(input GetPluginInfoInput) (dto.PluginInfoResult, error) {
	if input.PluginName == "" {
		return dto.PluginInfoResult{}, fmt.Errorf("plugin name is required")
	}

	// 1. Get plugin metadata
	plugin, err := uc.pluginRepo.Get(input.PluginName)
	if err != nil {
		return dto.PluginInfoResult{}, err
	}

	// 2. Read plugin state
	state, err := uc.stateStore.Read(input.PluginName)
	if err != nil {
		return dto.PluginInfoResult{}, err
	}

	// 3. Build Result
	result := dto.PluginInfoResult{
		Plugin: dto.PluginInfo{
			Name:        plugin.Name,
			Interpreter: plugin.Interpreter,
			Description: plugin.Description,
			Entry:       plugin.Entry,
			Usage:       plugin.Usage,
			Source:      plugin.Source,
			AddedAt:     plugin.AddedAt,
		},
		State: dto.PluginStateInfo{
			LastExecutedAt: state.LastExecutedAt,
			RunCount:       state.RunCount,
			MostRecentArgs: state.GetMostRecentArgs(),
		},
	}

	// 4. Extract history
	historyLimit := input.HistoryLimit
	if historyLimit <= 0 {
		historyLimit = 10 // Default
	}

	historyCount := len(state.History)
	start := historyCount - historyLimit
	if start < 0 {
		start = 0
	}

	for i := historyCount - 1; i >= start; i-- {
		record := state.History[i]
		result.RecentHistory = append(result.RecentHistory, dto.RunRecordInfo{
			Timestamp: record.Timestamp,
			Args:      record.Args,
			ExitCode:  record.ExitCode,
			Duration:  record.Duration,
			Status:    string(record.Status),
		})
	}

	return result, nil
}
