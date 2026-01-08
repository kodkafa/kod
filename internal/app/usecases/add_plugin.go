package usecases

import (
	"fmt"
	"net/url"
	"strings"

	"kodkafa/internal/app/dto"
	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
)

// AddPluginUseCase handles adding a new plugin.
type AddPluginUseCase struct {
	pluginRepo  ports.PluginRepository
	stateStore  ports.StateStore
	configStore ports.ConfigStore
}

// NewAddPluginUseCase creates a new AddPluginUseCase.
func NewAddPluginUseCase(pluginRepo ports.PluginRepository, stateStore ports.StateStore, configStore ports.ConfigStore) *AddPluginUseCase {
	return &AddPluginUseCase{
		pluginRepo:  pluginRepo,
		stateStore:  stateStore,
		configStore: configStore,
	}
}

// AddPluginInput represents the input for AddPluginUseCase.
type AddPluginInput struct {
	Source string
}

// Execute adds a plugin and initializes its state.
func (uc *AddPluginUseCase) Execute(input AddPluginInput) (dto.AddPluginResult, error) {
	if input.Source == "" {
		return dto.AddPluginResult{Success: false, Message: "source path is required"}, fmt.Errorf("source is required")
	}

	// 1. Check if source is a remote URL
	if strings.HasPrefix(input.Source, "http") || strings.Contains(input.Source, "@") {
		// Parse URL to check domain
		u, err := url.Parse(input.Source)
		if err != nil && !strings.Contains(input.Source, "@") {
			return dto.AddPluginResult{Success: false, Message: "invalid source URL"}, err
		}

		domain := ""
		if u != nil && u.Host != "" {
			domain = u.Host
		} else if strings.Contains(input.Source, "@") {
			// Handle git@github.com:...
			parts := strings.Split(strings.Split(input.Source, "@")[1], ":")
			domain = parts[0]
		}

		config, err := uc.configStore.Read()
		if err == nil { // Proceed with trusted domain check only if config can be read
			trusted := false
			for _, d := range config.TrustedDomains {
				if strings.Contains(domain, d) {
					trusted = true
					break
				}
			}
			if !trusted {
				return dto.AddPluginResult{Success: false, Message: "domain not trusted: " + domain}, fmt.Errorf("domain not trusted")
			}
		}
	}

	// 2. Add plugin via repository
	plugin, err := uc.pluginRepo.Add(input.Source)
	if err != nil {
		return dto.AddPluginResult{Success: false, Message: err.Error()}, err
	}

	// 3. Initialize state
	state := entities.NewPluginState(plugin.Name)
	if err := uc.stateStore.Write(state); err != nil {
		return dto.AddPluginResult{
			Success: true, // Plugin was added, but state failed
			Message: fmt.Sprintf("plugin added but failed to initialize state: %v", err),
			Plugin: dto.PluginInfo{
				Name:        plugin.Name,
				Interpreter: plugin.Interpreter,
				Description: plugin.Description,
				Entry:       plugin.Entry,
				Source:      plugin.Source,
				AddedAt:     plugin.AddedAt,
			},
		}, nil
	}

	return dto.AddPluginResult{
		Success: true,
		Message: "plugin added successfully",
		Plugin: dto.PluginInfo{
			Name:        plugin.Name,
			Interpreter: plugin.Interpreter,
			Description: plugin.Description,
			Entry:       plugin.Entry,
			Source:      plugin.Source,
			AddedAt:     plugin.AddedAt,
		},
	}, nil
}
