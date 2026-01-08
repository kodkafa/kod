package dto

import (
	"time"
)

// PluginListItem - for list view
type PluginListItem struct {
	Name        string    `json:"name"`
	Interpreter string    `json:"interpreter"`
	Description string    `json:"description"`
	Usage       string    `json:"usage"`
	LastRun     time.Time `json:"last_run"`  // from state
	RunCount    int       `json:"run_count"` // from state
}

// PluginInfo - detailed plugin info
type PluginInfo struct {
	Name        string    `json:"name"`
	Interpreter string    `json:"interpreter"`
	Description string    `json:"description"`
	Entry       string    `json:"entry"`
	Usage       string    `json:"usage"`
	Source      string    `json:"source"`
	AddedAt     time.Time `json:"added_at"`
}

// PluginStateInfo - state summary
type PluginStateInfo struct {
	LastExecutedAt time.Time `json:"last_executed_at"`
	RunCount       int       `json:"run_count"`
	MostRecentArgs string    `json:"most_recent_args"`
}

// RunRecordInfo - run record for display
type RunRecordInfo struct {
	Timestamp time.Time     `json:"timestamp"`
	Args      string        `json:"args"`
	ExitCode  int           `json:"exit_code"`
	Duration  time.Duration `json:"duration"`
	Status    string        `json:"status"`
}

// DashboardDTO - for ListPluginsUseCase
type DashboardDTO struct {
	TopPlugins  []PluginListItem `json:"top_plugins"`
	MainPlugins []PluginListItem `json:"main_plugins"`
	TotalCount  int              `json:"total_count"`
	CurrentPage int              `json:"current_page"`
	PageSize    int              `json:"page_size"`
	TotalPages  int              `json:"total_pages"`
	ShowTopList bool             `json:"show_top_list"`
}

// AddPluginResult - for AddPluginUseCase
type AddPluginResult struct {
	Plugin  PluginInfo `json:"plugin"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
}

// DeletePluginResult - for DeletePluginUseCase
type DeletePluginResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	PluginName string `json:"plugin_name"`
}

// RunPluginResult - for RunPluginUseCase
type RunPluginResult struct {
	PluginName string `json:"plugin_name"`
	Args       string `json:"args"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	// Note: result mapping from ports.RunResult can be added if needed
	ExitCode    int           `json:"exit_code"`
	Duration    time.Duration `json:"duration"`
	Status      string        `json:"status"`
	Interpreter string        `json:"interpreter"`
	Output      string        `json:"output"`
}

// PluginInfoResult - for GetPluginInfoUseCase
type PluginInfoResult struct {
	Plugin        PluginInfo      `json:"plugin"`
	State         PluginStateInfo `json:"state"`
	RecentHistory []RunRecordInfo `json:"recent_history"`
}
