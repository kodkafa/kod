package usecases

import (
	"math"
	"sort"

	"kodkafa/internal/app/dto"
	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
)

// ListPluginsUseCase handles listing plugins with a dual-list dashboard structure.
type ListPluginsUseCase struct {
	pluginRepo  ports.PluginRepository
	usageStore  ports.UsageStore
	configStore ports.ConfigStore
	stateStore  ports.StateStore
}

// NewListPluginsUseCase creates a new ListPluginsUseCase.
func NewListPluginsUseCase(
	pluginRepo ports.PluginRepository,
	usageStore ports.UsageStore,
	configStore ports.ConfigStore,
	stateStore ports.StateStore,
) *ListPluginsUseCase {
	return &ListPluginsUseCase{
		pluginRepo:  pluginRepo,
		usageStore:  usageStore,
		configStore: configStore,
		stateStore:  stateStore,
	}
}

// ListPluginsInput represents the input for the ListPluginsUseCase.
type ListPluginsInput struct {
	Page     int
	PageSize int
}

// Execute returns a DashboardDTO with top and main lists.
func (uc *ListPluginsUseCase) Execute(input ListPluginsInput) (dto.DashboardDTO, error) {
	result := dto.DashboardDTO{
		CurrentPage: input.Page,
		PageSize:    input.PageSize,
	}

	// 1. Read config
	config, err := uc.configStore.Read()
	if err != nil {
		return result, err
	}

	// Determine page size
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = config.ItemsPerPage
		if pageSize <= 0 {
			pageSize = 10 // Default fallback
		}
	}
	result.PageSize = pageSize

	// Check preferences for showing top list
	if !config.ShowLastRuns {
		result.ShowTopList = false
		// If top list is not shown, we can skip fetching usage stats and building top list
		// and proceed directly to fetching all plugins and building the main list.
		// However, the original logic fetches all plugins first, so we'll keep that order.
	} else {
		result.ShowTopList = true
	}

	// 2. Fetch all plugins
	plugins, err := uc.pluginRepo.List()
	if err != nil {
		return result, err
	}

	// 3. Read usage stats
	usage, err := uc.usageStore.Read()
	if err != nil {
		// If usage stats don't exist, we continue with empty lists
		usage = nil
	}

	// Prepare state map for efficient lookup
	stateMap := make(map[string]*entities.PluginState)
	for _, pl := range plugins {
		state, _ := uc.stateStore.Read(pl.Name)
		if state != nil {
			stateMap[pl.Name] = state
		}
	}

	// 4. Build Top List
	topPluginsMap := make(map[string]bool)
	if config.ShowLastRuns && usage != nil {
		var topItems []string
		if config.LastRunOrder == "most" {
			for i := 0; i < len(usage.MostUsed) && i < config.FavLimit; i++ {
				topItems = append(topItems, usage.MostUsed[i].PluginName)
			}
		} else { // default "last"
			for i := 0; i < len(usage.RecentlyUsed) && i < config.FavLimit; i++ {
				topItems = append(topItems, usage.RecentlyUsed[i].PluginName)
			}
		}

		for _, name := range topItems {
			// Find plugin metadata
			var p *dto.PluginListItem
			for _, pl := range plugins {
				if pl.Name == name {
					state, _ := uc.stateStore.Read(pl.Name)
					p = &dto.PluginListItem{
						Name:        pl.Name,
						Interpreter: pl.Interpreter,
						Description: pl.Description,
					}
					if state != nil {
						p.LastRun = state.LastExecutedAt
						p.RunCount = state.RunCount
					}
					break
				}
			}
			if p != nil {
				result.TopPlugins = append(result.TopPlugins, *p)
				topPluginsMap[name] = true
			}
		}
	}

	// 5. Build Main List
	var mainList []dto.PluginListItem
	for _, pl := range plugins {
		if !topPluginsMap[pl.Name] {
			state, _ := uc.stateStore.Read(pl.Name)
			item := dto.PluginListItem{
				Name:        pl.Name,
				Interpreter: pl.Interpreter,
				Description: pl.Description,
			}
			if state != nil {
				item.LastRun = state.LastExecutedAt
				item.RunCount = state.RunCount
			}
			mainList = append(mainList, item)
		}
	}

	// Sort main list alphabetically
	sort.Slice(mainList, func(i, j int) bool {
		return mainList[i].Name < mainList[j].Name
	})

	result.TotalCount = len(mainList)
	result.TotalPages = int(math.Ceil(float64(result.TotalCount) / float64(pageSize)))
	if result.TotalPages == 0 {
		result.TotalPages = 1
	}

	// Apply pagination
	start := (input.Page - 1) * pageSize
	if start < 0 {
		start = 0
	}
	end := start + pageSize
	if start > len(mainList) {
		result.MainPlugins = []dto.PluginListItem{}
	} else {
		if end > len(mainList) {
			end = len(mainList)
		}
		result.MainPlugins = mainList[start:end]
	}

	return result, nil
}
