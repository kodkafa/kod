package usecases

import (
	"fmt"
	"time"

	"kodkafa/internal/app/dto"
	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
)

// RunPluginUseCase handles the lifecycle of executing a plugin.
type RunPluginUseCase struct {
	pluginRepo  ports.PluginRepository
	stateStore  ports.StateStore
	usageStore  ports.UsageStore
	configStore ports.ConfigStore
	runner      ports.Runner
}

// NewRunPluginUseCase creates a new RunPluginUseCase.
func NewRunPluginUseCase(
	pluginRepo ports.PluginRepository,
	stateStore ports.StateStore,
	usageStore ports.UsageStore,
	configStore ports.ConfigStore,
	runner ports.Runner,
) *RunPluginUseCase {
	return &RunPluginUseCase{
		pluginRepo:  pluginRepo,
		stateStore:  stateStore,
		usageStore:  usageStore,
		configStore: configStore,
		runner:      runner,
	}
}

// RunPluginInput represents the input for RunPluginUseCase.
type RunPluginInput struct {
	PluginName string
	Args       string
	Mode       ports.RunMode
	OutputChan chan<- ports.OutputChunk
}

// Execute orchestrates the plugin execution lifecycle.
func (uc *RunPluginUseCase) Execute(input RunPluginInput) (dto.RunPluginResult, error) {
	if input.PluginName == "" {
		return dto.RunPluginResult{Success: false, Status: "error"}, fmt.Errorf("plugin name is required")
	}

	result := dto.RunPluginResult{
		PluginName: input.PluginName,
		Args:       input.Args,
	}

	// 1. Get plugin
	plugin, err := uc.pluginRepo.Get(input.PluginName)
	if err != nil {
		result.Status = "not_found"
		return result, err
	}

	// Set interpreter for display
	switch plugin.Interpreter {
	case "python":
		result.Interpreter = "python3"
	case "node":
		result.Interpreter = "node"
	default:
		result.Interpreter = plugin.Interpreter
	}

	// 2. Update usage stats (Run started - P0)
	usage, err := uc.usageStore.Read()
	if err == nil && usage != nil {
		config, err := uc.configStore.Read()
		limit := 10 // Default
		if err == nil {
			limit = config.LastRunLimit
		}
		// Update global usage (for favorites) - no history here
		usage.RecordRun(input.PluginName, limit)
		_ = uc.usageStore.Write(usage)
	}

	// 3. Update plugin state (Persist run intent - P0)
	state, err := uc.stateStore.Read(input.PluginName)
	if err != nil {
		state = entities.NewPluginState(input.PluginName)
	}

	record := entities.RunRecord{
		Timestamp: time.Now(),
		Args:      input.Args, // History lives here!
		Status:    entities.RunStatusRunning,
	}
	state.AddRunRecord(record)
	_ = uc.stateStore.Write(state)

	// 4. Run plugin (P1/P1i)
	runResult, err := uc.runner.Run(plugin, input.Args, input.Mode, input.OutputChan)

	// 5. Finalize record (P2/P3)
	if err != nil {
		record.Status = entities.RunStatusFailed
		result.Status = string(entities.RunStatusFailed)
		result.Success = false
		result.Message = err.Error()
	} else {
		record.ExitCode = runResult.ExitCode
		record.Duration = time.Duration(runResult.Duration)
		record.Status = entities.RunStatus(runResult.Status)

		result.ExitCode = runResult.ExitCode
		result.Duration = time.Duration(runResult.Duration)
		result.Status = runResult.Status
		result.Output = runResult.Output
		result.Success = runResult.ExitCode == 0

		if runResult.ExitCode != 0 {
			result.Message = fmt.Sprintf("Process exited with code %d", runResult.ExitCode)
		}
	}

	// Update the record with final results
	if len(state.History) > 0 {
		state.History[len(state.History)-1] = record
		_ = uc.stateStore.Write(state)
	}

	return result, err
}
