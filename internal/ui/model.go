package ui

import (
	"fmt"
	"kodkafa/internal/app/dto"
	"kodkafa/internal/app/usecases"
	"kodkafa/internal/domain/ports"
	"kodkafa/internal/ui/screens"
	"kodkafa/internal/ui/tea"

	tea_pkg "github.com/charmbracelet/bubbletea"
)

// Model is the root TUI model.
type Model struct {
	state        tea.TUIState
	dashboard    tea_pkg.Model
	commandMenu  tea_pkg.Model
	inputModel   tea_pkg.Model
	promptModel  tea_pkg.Model
	runningModel tea_pkg.Model
	postRunModel tea_pkg.Model
	infoModel    tea_pkg.Model
	activeScreen tea_pkg.Model
	initialCmd   tea_pkg.Cmd

	outputChan chan ports.OutputChunk

	isLoading bool

	listUC   *usecases.ListPluginsUseCase
	addUC    *usecases.AddPluginUseCase
	deleteUC *usecases.DeletePluginUseCase
	loadUC   *usecases.LoadPluginDepsUseCase
	infoUC   *usecases.GetPluginInfoUseCase
	runUC    *usecases.RunPluginUseCase
	initUC   *usecases.InitLayoutUseCase

	pendingCmd        string
	pendingName       string
	deletePendingName string
	deleteRemoveDeps  bool
}

// NewModel creates the root TUI model.
func NewModel(listUC *usecases.ListPluginsUseCase, addUC *usecases.AddPluginUseCase, deleteUC *usecases.DeletePluginUseCase, loadUC *usecases.LoadPluginDepsUseCase, infoUC *usecases.GetPluginInfoUseCase, runUC *usecases.RunPluginUseCase, initUC *usecases.InitLayoutUseCase, showSplash bool) *Model {
	dashboard := screens.NewDashboardModel(listUC)

	var activeScreen tea_pkg.Model = dashboard
	if showSplash {
		activeScreen = screens.NewSplashModel()
	}

	return &Model{
		state:        tea.StateNormal,
		activeScreen: activeScreen,
		dashboard:    dashboard,
		listUC:       listUC,
		addUC:        addUC,
		deleteUC:     deleteUC,
		loadUC:       loadUC,
		infoUC:       infoUC,
		runUC:        runUC,
		initUC:       initUC,
	}
}

// StartRun configures the model to start immediately with the run prompt for the given plugin.
func (m *Model) StartRun(pluginName string) {
	m.pendingCmd = "kod info"
	m.initialCmd = func() tea_pkg.Msg {
		return tea.PluginSelectedMsg{PluginName: pluginName}
	}
}

// Init initializes the root model.
func (m *Model) Init() tea_pkg.Cmd {
	if m.initialCmd != nil {
		return tea_pkg.Batch(m.activeScreen.Init(), m.initialCmd)
	}
	return m.activeScreen.Init()
}

// Update handles top-level messages and delegates to screens.
func (m *Model) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	var cmd tea_pkg.Cmd

	switch msg := msg.(type) {
	case tea.SwitchScreenMsg:
		if msg.ScreenName == "dashboard" {
			m.activeScreen = m.dashboard
			return m, m.activeScreen.Init()
		}

	case tea.PluginSelectedMsg:
		m.pendingName = msg.PluginName
		cmdToRun := msg.Cmd
		if cmdToRun == "" {
			cmdToRun = m.pendingCmd
		}
		if cmdToRun == "" {
			cmdToRun = "kod run"
		}
		m.pendingCmd = cmdToRun

		switch cmdToRun {
		case "kod del":
			m.deletePendingName = msg.PluginName
			m.state = tea.StateDeleteConfirm
			m.activeScreen = screens.NewConfirmModel("DELETE PLUGIN", msg.PluginName, "This will remove the plugin source and metadata.", func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateDeleteDepsConfirm}
			}, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateCommandMenu}
			})
			return m, m.activeScreen.Init()
		case "kod load":
			return m, func() tea_pkg.Msg {
				return tea.LoadPluginMsg{PluginName: msg.PluginName}
			}
		case "kod run":
			m.loading(true)
			return m, func() tea_pkg.Msg {
				res, err := m.infoUC.Execute(usecases.GetPluginInfoInput{PluginName: msg.PluginName})
				if err != nil {
					return tea.ErrMsg{Err: err}
				}
				return tea.PluginInfoFetchedMsg{Data: res}
			}
		case "kod info":
			m.loading(true)
			return m, func() tea_pkg.Msg {
				res, err := m.infoUC.Execute(usecases.GetPluginInfoInput{PluginName: msg.PluginName})
				if err != nil {
					return tea.ErrMsg{Err: err}
				}
				return tea.PluginInfoFetchedMsg{Data: res}
			}
		case "kod log":
			// TODO: Implement log view
			return m, nil
		}

	case tea.PluginInfoFetchedMsg:
		m.loading(false)

		// Decide where to go based on pendingCmd
		if m.pendingCmd == "kod info" {
			// Show Info Screen
			m.infoModel = screens.NewInfoModel(msg.Data)
			m.activeScreen = m.infoModel
			// Info screen has no particular state enum, we can reuse StateNormal or add StateInfo
			// Let's keep StateNormal effectively or just not change state enum but activeScreen matches.
			// Actually, to support Back properly, let's treat it as a screen.
			return m, m.infoModel.Init()
		}

		// Default or "kod run" -> Show Prompt
		m.state = tea.StatePrompt

		// Extract argument history for prompt navigation
		var argsHistory []string
		for _, h := range msg.Data.RecentHistory {
			if h.Args != "" {
				argsHistory = append(argsHistory, h.Args)
			}
		}

		m.promptModel = screens.NewPromptModel(msg.Data.Plugin, argsHistory, m.runUC)
		m.activeScreen = m.promptModel
		return m, m.promptModel.Init()

	case tea.PluginRunMsg:
		m.state = tea.StateRunning
		m.runningModel = screens.NewRunningModel(msg.PluginName)
		m.activeScreen = m.runningModel

		m.outputChan = make(chan ports.OutputChunk)

		// 1. Command to start execution
		runCmd := func() tea_pkg.Msg {
			res, _ := m.runUC.Execute(usecases.RunPluginInput{
				PluginName: msg.PluginName,
				Args:       msg.Args,
				Mode:       ports.RunModeStreaming,
				OutputChan: m.outputChan,
			})
			return tea.RunFinishedMsg{Result: res}
		}

		return m, tea_pkg.Batch(m.activeScreen.Init(), runCmd, waitForOutput(m.outputChan))

	case tea.OutputMsg:
		var innerCmd tea_pkg.Cmd
		m.runningModel, innerCmd = m.runningModel.Update(msg)
		return m, tea_pkg.Batch(innerCmd, waitForOutput(m.outputChan))

	case tea.LoadPluginMsg:
		m.loading(true)
		return m, func() tea_pkg.Msg {
			_, err := m.loadUC.Execute(usecases.LoadPluginDepsInput{PluginName: msg.PluginName})
			if err != nil {
				return tea.ErrMsg{Err: err}
			}
			return tea.PluginLoadedMsg{PluginName: msg.PluginName}
		}

	case tea.DeletePluginMsg:
		m.deletePendingName = msg.PluginName
		m.deleteRemoveDeps = msg.RemoveDeps
		m.loading(true)
		return m, m.deleteCmd()

	case tea.InitLayoutMsg:
		m.loading(true)
		return m, func() tea_pkg.Msg {
			err := m.initUC.Execute()
			res := dto.RunPluginResult{
				PluginName: "System Init",
				Success:    err == nil,
				Status:     "complete",
			}
			if err != nil {
				res.Status = "failed"
				return tea.ErrMsg{Err: err}
			}
			return tea.RunFinishedMsg{Result: res}
		}

	case tea.PluginLoadedMsg:
		m.loading(false)
		m.state = tea.StateNormal
		m.activeScreen = m.dashboard
		// Refresh dashboard to show new plugin or status
		return m, m.dashboard.Init()

	case tea.RunFinishedMsg:
		m.loading(false)
		m.state = tea.StatePostRun
		m.postRunModel = screens.NewResultsModel(msg.Result)
		m.activeScreen = m.postRunModel
		return m, m.postRunModel.Init()

	case tea.SwitchStateMsg:
		m.state = msg.State
		switch m.state {
		case tea.StateNormal:
			m.activeScreen = m.dashboard
			return m, m.dashboard.Init()
		case tea.StateCommandMenu:
			m.commandMenu = screens.NewMenuModel()
			m.activeScreen = m.commandMenu
			m.pendingCmd = msg.Cmd
			return m, m.commandMenu.Init()
		case tea.StateInput:
			m.pendingCmd = msg.Cmd
			m.inputModel = screens.NewInputModel(m.addUC, msg.Mode)
			m.activeScreen = m.inputModel
			return m, m.inputModel.Init()
		case tea.StateDeleteConfirm:
			// Handled by switch cmdToRun
		case tea.StateDeleteDepsConfirm:
			m.activeScreen = screens.NewConfirmModel("REMOVE DEPENDENCIES?", m.deletePendingName, "Dependencies will be removed. Run 'kod load' to reinstall.", func() tea_pkg.Msg {
				m.deleteRemoveDeps = true
				return m.deleteCmd()()
			}, func() tea_pkg.Msg {
				m.deleteRemoveDeps = false
				return m.deleteCmd()()
			})
			return m, m.activeScreen.Init()
		case tea.StatePrompt:
			if m.promptModel == nil {
				// If we don't have a prompt model (e.g. from 'init' command),
				// fallback to dashboard instead of crashing.
				m.state = tea.StateNormal
				m.activeScreen = m.dashboard
				return m, m.dashboard.Init()
			}
			m.activeScreen = m.promptModel
			return m, m.promptModel.Init()
		}

	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == tea.StateNormal {
				return m, tea_pkg.Quit
			}
		}
	case tea.PluginAddedMsg:
		m.loading(true)
		// Auto-load dependencies
		return m, func() tea_pkg.Msg {
			_, err := m.loadUC.Execute(usecases.LoadPluginDepsInput{PluginName: msg.PluginName})
			if err != nil {
				return tea.ErrMsg{Err: err}
			}
			return tea.PluginLoadedMsg{PluginName: msg.PluginName}
		}

	case tea.ErrMsg:
		m.loading(false)
	}

	m.activeScreen, cmd = m.activeScreen.Update(msg)
	return m, cmd
}

func (m *Model) loading(b bool) {
	m.isLoading = b
}

func (m *Model) deleteCmd() tea_pkg.Cmd {
	return func() tea_pkg.Msg {
		_, err := m.deleteUC.Execute(usecases.DeletePluginInput{
			PluginName: m.deletePendingName,
			RemoveDeps: m.deleteRemoveDeps,
		})
		if err != nil {
			return tea.ErrMsg{Err: err}
		}
		return tea.SwitchStateMsg{State: tea.StateNormal}
	}
}

func waitForOutput(c chan ports.OutputChunk) tea_pkg.Cmd {
	return func() tea_pkg.Msg {
		chunk, ok := <-c
		if !ok {
			return nil
		}
		return tea.OutputMsg{Chunk: string(chunk.Data)}
	}
}

// View renders the current active screen.
func (m *Model) View() string {
	view := m.activeScreen.View()
	if m.isLoading {
		return fmt.Sprintf("\n  Thinking...\n\n%s", view) // Overlay or prefix
	}
	return view
}
