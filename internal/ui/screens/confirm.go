package screens

import (
	"fmt"

	"kodkafa/internal/ui/theme"

	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmModel struct {
	title      string
	pluginName string
	note       string
	onConfirm  func() tea_pkg.Msg
	onCancel   func() tea_pkg.Msg
}

func NewConfirmModel(title, pluginName, note string, onConfirm, onCancel func() tea_pkg.Msg) *ConfirmModel {
	return &ConfirmModel{
		title:      title,
		pluginName: pluginName,
		note:       note,
		onConfirm:  onConfirm,
		onCancel:   onCancel,
	}
}

func (m *ConfirmModel) Init() tea_pkg.Cmd {
	return nil
}

func (m *ConfirmModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	switch msg := msg.(type) {
	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			return m, m.onConfirm
		case "n", "N", "esc":
			return m, m.onCancel
		}
	}
	return m, nil
}

func (m *ConfirmModel) View() string {
	style := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	noteStr := ""
	if m.note != "" {
		noteStr = lipgloss.NewStyle().Foreground(theme.TextSecondary).Render(fmt.Sprintf("\n  Note: %s", m.note))
	}

	return fmt.Sprintf(
		"\n  %s: %s%s\n\n  Are you sure? (y/N)",
		m.title,
		style.Render(m.pluginName),
		noteStr,
	)
}
