package screens

import (
	"fmt"
	"kodkafa/internal/app/dto"
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"
	"strings"

	"kodkafa/internal/ui/theme"

	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(theme.Success)
	failStyle    = lipgloss.NewStyle().Foreground(theme.Error)
)

type ResultsModel struct {
	result dto.RunPluginResult
}

func NewResultsModel(res dto.RunPluginResult) *ResultsModel {
	return &ResultsModel{result: res}
}

func (m *ResultsModel) Init() tea_pkg.Cmd {
	return nil
}

func (m *ResultsModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	switch msg := msg.(type) {
	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "enter":
			// Return to Prompt - Enter
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StatePrompt}
			}
		case "esc":
			// Back - A
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateNormal}
			}
		}
	}
	return m, nil
}

func (m *ResultsModel) View() string {
	var b strings.Builder
	b.WriteString(components.RenderHeader("", "RESULTS"))

	status := successStyle.Render("SUCCESS")
	if !m.result.Success {
		status = failStyle.Render("FAILED")
		if m.result.Message != "" {
			status += "\n  Error:       " + failStyle.Render(m.result.Message)
		}
	}

	// PROMPT: blue(interpreter) yellow(plugin) args
	interpreter := lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true).Render(m.result.Interpreter)
	plugin := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render(m.result.PluginName)

	b.WriteString(fmt.Sprintf("\nPROMPT: %s %s %s\n\n", interpreter, plugin, m.result.Args))

	b.WriteString(fmt.Sprintf("  Status:      %s\n", status))
	b.WriteString(fmt.Sprintf("  ExitCode:    %d\n", m.result.ExitCode))
	b.WriteString(fmt.Sprintf("  Duration:    %v\n", m.result.Duration))

	if m.result.Output != "" {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(theme.Muted).Render("OUTPUT"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(theme.Muted).PaddingLeft(1).Render(m.result.Output))
	}

	b.WriteString("\n\n")
	b.WriteString(components.RenderFooter(
		components.FooterItem{Key: "Esc", Label: "Back"},
		components.FooterItem{Key: "Enter", Label: "Re-run"},
	))

	return b.String()
}
