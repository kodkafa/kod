package usecases

import (
	"fmt"
	"kodkafa/internal/app/dto"
	"kodkafa/internal/domain/ports"
)

type LoadPluginDepsUseCase struct {
	pluginRepo ports.PluginRepository
	installer  ports.DependencyInstaller
}

func NewLoadPluginDepsUseCase(pluginRepo ports.PluginRepository, installer ports.DependencyInstaller) *LoadPluginDepsUseCase {
	return &LoadPluginDepsUseCase{
		pluginRepo: pluginRepo,
		installer:  installer,
	}
}

type LoadPluginDepsInput struct {
	PluginName string
}

func (uc *LoadPluginDepsUseCase) Execute(input LoadPluginDepsInput) (dto.RunPluginResult, error) {
	result := dto.RunPluginResult{
		PluginName: input.PluginName,
	}

	plugin, err := uc.pluginRepo.Get(input.PluginName)
	if err != nil {
		result.Message = "plugin not found"
		result.Success = false
		return result, err
	}

	err = uc.installer.Install(plugin)
	if err != nil {
		result.Message = fmt.Sprintf("failed to install dependencies: %v", err)
		result.Success = false
		return result, err
	}

	result.Message = "Dependencies installed successfully"
	result.Status = "completed"
	result.Success = true
	return result, nil
}
