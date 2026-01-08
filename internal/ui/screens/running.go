package screens

import (
	"fmt"
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"
	"strings"

	"kodkafa/internal/ui/theme"

	"github.com/charmbracelet/bubbles/spinner"
	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RunningModel struct {
	pluginName string
	spinner    spinner.Model
	logs       []string // Added
}

func NewRunningModel(pluginName string) *RunningModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.Primary)
	return &RunningModel{
		pluginName: pluginName,
		spinner:    s,
		logs:       make([]string, 0), // Added
	}
}

func (m *RunningModel) Init() tea_pkg.Cmd {
	return m.spinner.Tick
}

func (m *RunningModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	var cmd tea_pkg.Cmd
	switch msg := msg.(type) {
	case tea.OutputMsg: // Added
		m.logs = append(m.logs, msg.Chunk)
		if len(m.logs) > 10 { // Keep last 10 lines
			m.logs = m.logs[len(m.logs)-10:]
		}
		return m, nil
	case tea.RunFinishedMsg:
		return m, func() tea_pkg.Msg {
			return tea.SwitchStateMsg{State: tea.StatePostRun} // Modified
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *RunningModel) View() string {
	var b strings.Builder

	b.WriteString(components.RenderHeader("A Persistent CLI with Memory", "RUN"))
	b.WriteString(fmt.Sprintf("\n  %s Running %s...\n\n", m.spinner.View(), m.pluginName)) // Modified

	if len(m.logs) > 0 { // Added
		b.WriteString(lipgloss.NewStyle().Foreground(theme.TextSecondary).Render("Output:") + "\n") // Added
		for _, line := range m.logs {                                                               // Added
			b.WriteString("  " + line + "\n") // Added
		}
	}

	return b.String() // Modified
}
