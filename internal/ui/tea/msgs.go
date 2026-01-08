package tea

import (
	"kodkafa/internal/app/dto"
)

// TickMsg is sent for animation ticks
type TickMsg struct{}

// PluginsLoadedMsg is sent when the dashboard data is ready
type PluginsLoadedMsg struct {
	Data dto.DashboardDTO
}

// PluginSelectedMsg is sent when a plugin is chosen from the dashboard
type PluginSelectedMsg struct {
	PluginName string
	Cmd        string
}

// PluginAddedMsg is sent when a plugin is successfully added
type PluginAddedMsg struct {
	PluginName string
}

// PluginInfoFetchedMsg is sent when detailed plugin info is loaded
type PluginInfoFetchedMsg struct {
	Data dto.PluginInfoResult
}

// PluginRunMsg is sent to start a plugin execution
type PluginRunMsg struct {
	PluginName string
	Args       string
}

// OutputMsg is sent when a plugin produces output
type OutputMsg struct {
	Chunk string
}

// RunFinishedMsg is sent when a plugin execution completes
type RunFinishedMsg struct {
	Result dto.RunPluginResult
}

// SwitchScreenMsg is sent to change the active TUI screen
type SwitchScreenMsg struct {
	ScreenName string
}

// SwitchStateMsg is sent to change the operational state of the TUI
type SwitchStateMsg struct {
	State TUIState
	Mode  InputModeType
	Cmd   string
}

// ErrMsg is sent when an application error occurs
type ErrMsg struct {
	Err error
}

type DeletePluginMsg struct {
	PluginName string
	RemoveDeps bool
}

type LoadPluginMsg struct {
	PluginName string
}

type InitLayoutMsg struct{}

// PluginLoadedMsg is sent when a plugin's dependencies are loaded efficiently
type PluginLoadedMsg struct {
	PluginName string
}
