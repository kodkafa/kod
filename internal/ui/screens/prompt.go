package screens

import (
	"fmt"
	"kodkafa/internal/app/dto"
	"kodkafa/internal/app/usecases"
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"
	"strings"

	"kodkafa/internal/ui/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.Accent).
				MarginBottom(1)

	usageStyle = lipgloss.NewStyle().
			Foreground(theme.Primary).
			Background(theme.FooterBg).
			Padding(0, 1)

	historyLabelStyle = lipgloss.NewStyle().
				Foreground(theme.TextSecondary).
				MarginTop(1)
)

type PromptModel struct {
	textInput     textinput.Model
	runUC         *usecases.RunPluginUseCase
	pluginInfo    dto.PluginInfo
	history       []string
	historyCursor int
	loading       bool
	err           error
}

func NewPromptModel(pluginInfo dto.PluginInfo, history []string, runUC *usecases.RunPluginUseCase) *PromptModel {
	ti := textinput.New()
	ti.Placeholder = "arguments..."
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 60

	// Start EMPTY, allow Up arrow to fetch history
	// if len(history) > 0 {
	// 	ti.SetValue(history[0])
	// }

	return &PromptModel{
		textInput:     ti,
		pluginInfo:    pluginInfo,
		history:       history,
		historyCursor: -1, // Start before first item
		runUC:         runUC,
	}
}

func (m *PromptModel) Init() tea_pkg.Cmd {
	m.loading = false
	return textinput.Blink
}

func (m *PromptModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	var cmd tea_pkg.Cmd

	switch msg := msg.(type) {
	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateNormal}
			}
		case "enter":
			args := m.textInput.Value()
			m.loading = true
			return m, func() tea_pkg.Msg {
				return tea.PluginRunMsg{PluginName: m.pluginInfo.Name, Args: args}
			}
		case "up":
			// Go to older history (increment cursor index in our list 0..N)
			if len(m.history) > 0 {
				if m.historyCursor < len(m.history)-1 {
					m.historyCursor++
					m.textInput.SetValue(m.history[m.historyCursor])
				}
			}
		case "down":
			// Go to newer (decrement cursor)
			if m.historyCursor > 0 {
				m.historyCursor--
				m.textInput.SetValue(m.history[m.historyCursor])
			} else if m.historyCursor == 0 {
				// Back to empty / current
				m.historyCursor = -1
				m.textInput.SetValue("")
			}
		}

	case tea.ErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *PromptModel) View() string {
	if m.loading {
		return "\n  Preparing execution..."
	}

	var b strings.Builder

	// Header
	b.WriteString(components.RenderHeader("A Persistent CLI with Memory", fmt.Sprintf("RUN PLUGIN: %s", m.pluginInfo.Name)))

	// Usage
	if m.pluginInfo.Usage != "" {
		b.WriteString("Usage:\n")
		b.WriteString(usageStyle.Render(m.pluginInfo.Usage) + "\n\n")
	}

	interpreter := lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true).Render(m.pluginInfo.Interpreter)
	plugin := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render(m.pluginInfo.Name)

	b.WriteString(fmt.Sprintf("%s %s ", interpreter, plugin))
	b.WriteString(m.textInput.View() + "\n")

	if len(m.history) > 0 {
		b.WriteString(historyLabelStyle.Render(fmt.Sprintf("History: %d/%d (Use ↑/↓ to navigate)", m.historyCursor+1, len(m.history))) + "\n")
	}

	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(theme.Error).Render(fmt.Sprintf("\nError: %v", m.err)))
	}

	b.WriteString("\n")
	b.WriteString(components.RenderFooter(
		components.FooterItem{Key: "↑/↓", Label: "History"},
		components.FooterItem{Key: "Esc", Label: "Cancel"},
		components.FooterItem{Key: "Enter", Label: "Run"},
	))

	return b.String()
}
