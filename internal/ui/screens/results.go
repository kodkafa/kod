package screens

import (
	"fmt"
	"strings"

	"kodkafa/internal/app/dto"
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"
	"kodkafa/internal/ui/theme"

	"github.com/charmbracelet/bubbles/viewport"
	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(theme.Success)
	failStyle    = lipgloss.NewStyle().Foreground(theme.Error)
)

type ResultsModel struct {
	result   dto.RunPluginResult
	viewport viewport.Model
	ready    bool
}

func NewResultsModel(res dto.RunPluginResult, width, height int) *ResultsModel {
	m := &ResultsModel{
		result: res,
	}
	if width > 0 && height > 0 {
		m.updateViewport(width, height)
	}
	return m
}

func (m *ResultsModel) updateViewport(width, height int) {
	headerHeight := 8 // Approximate height of header + status info
	footerHeight := 3 // Approximate height of footer
	verticalMarginHeight := headerHeight + footerHeight

	m.viewport = viewport.New(width, height-verticalMarginHeight)
	m.viewport.YPosition = headerHeight
	m.viewport.SetContent(m.result.Output)
	m.ready = true
}

func (m *ResultsModel) Init() tea_pkg.Cmd {
	return nil
}

func (m *ResultsModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	var (
		cmd  tea_pkg.Cmd
		cmds []tea_pkg.Cmd
	)

	switch msg := msg.(type) {
	case tea_pkg.WindowSizeMsg:
		m.updateViewport(msg.Width, msg.Height)

	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StatePrompt}
			}
		case "esc":
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateNormal}
			}
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea_pkg.Batch(cmds...)
}

func (m *ResultsModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var b strings.Builder

	// Create header content
	b.WriteString(components.RenderHeader("", "RESULTS"))

	status := successStyle.Render("SUCCESS")
	if !m.result.Success {
		status = failStyle.Render("FAILED")
		if m.result.Message != "" {
			status += "\n  Error:       " + failStyle.Render(m.result.Message)
		}
	}

	interpreter := lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true).Render(m.result.Interpreter)
	plugin := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render(m.result.PluginName)

	b.WriteString(fmt.Sprintf("\nPROMPT: %s %s %s\n\n", interpreter, plugin, m.result.Args))
	b.WriteString(fmt.Sprintf("  Status:      %s\n", status))
	b.WriteString(fmt.Sprintf("  ExitCode:    %d\n", m.result.ExitCode))
	b.WriteString(fmt.Sprintf("  Duration:    %v\n", m.result.Duration))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.Muted).Render("OUTPUT (Scroll with ↑/↓)"))
	b.WriteString("\n")

	// Render viewport
	b.WriteString(m.viewport.View())

	// Ensure we don't double newlines too much, viewport might handle it.

	b.WriteString("\n")
	b.WriteString(components.RenderFooter(
		components.FooterItem{Key: "Esc", Label: "Back"},
		components.FooterItem{Key: "Enter", Label: "Re-run"},
		components.FooterItem{Key: "↑/↓", Label: "Scroll"},
	))

	return b.String()
}
