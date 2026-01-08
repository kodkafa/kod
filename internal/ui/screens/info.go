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
	infoLabelStyle = lipgloss.NewStyle().Foreground(theme.TextSecondary).Width(12)
	infoValueStyle = lipgloss.NewStyle().Foreground(theme.TextPrimary)
	primaryStyle   = lipgloss.NewStyle().Foreground(theme.Primary)
	secondaryStyle = lipgloss.NewStyle().Foreground(theme.Secondary)
)

type InfoModel struct {
	data dto.PluginInfoResult
}

func NewInfoModel(data dto.PluginInfoResult) *InfoModel {
	return &InfoModel{data: data}
}

func (m *InfoModel) Init() tea_pkg.Cmd {
	return nil
}

func (m *InfoModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	switch msg := msg.(type) {
	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "esc", "left":
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateNormal}
			}
		case "enter":
			// Transition to Run (via Prompt)
			return m, func() tea_pkg.Msg {
				return tea.PluginSelectedMsg{PluginName: m.data.Plugin.Name, Cmd: "kod run"}
			}
		}
	}
	return m, nil
}

func (m *InfoModel) View() string {
	var b strings.Builder

	b.WriteString(components.RenderHeader("A Persistent CLI with Memory", "PLUGIN INFO"))

	p := m.data.Plugin
	s := m.data.State

	// Metadata
	b.WriteString(fmt.Sprintf("%s %s\n", infoLabelStyle.Render("Name:"), infoValueStyle.Render(p.Name)))
	b.WriteString(fmt.Sprintf("%s %s\n", infoLabelStyle.Render("Interpreter:"), infoValueStyle.Render(p.Interpreter)))
	b.WriteString(fmt.Sprintf("%s %s\n", infoLabelStyle.Render("Description:"), infoValueStyle.Render(p.Description)))
	b.WriteString(fmt.Sprintf("%s %s\n", infoLabelStyle.Render("Added:"), infoValueStyle.Render(p.AddedAt.Format("2006-01-02 15:04"))))
	b.WriteString(fmt.Sprintf("%s %s\n", infoLabelStyle.Render("Source:"), infoValueStyle.Render(p.Source)))
	b.WriteString(fmt.Sprintf("%s %s\n", primaryStyle.Render("Usage:"), secondaryStyle.Render(p.Usage)))

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(theme.Accent).Render("STATS"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%s %s\n", infoLabelStyle.Render("Run Count:"), infoValueStyle.Render(fmt.Sprintf("%d", s.RunCount))))
	b.WriteString(fmt.Sprintf("%s %s\n", infoLabelStyle.Render("Last Run:"), infoValueStyle.Render(s.LastExecutedAt.Format("2006-01-02 15:04"))))

	if len(m.data.RecentHistory) > 0 {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(theme.Accent).Render("RECENT HISTORY"))
		b.WriteString("\n")
		for i, h := range m.data.RecentHistory {
			if i >= 5 {
				break
			} // Limit to 5
			status := "✓"
			if h.Status != "completed" {
				status = "✗"
			}
			b.WriteString(fmt.Sprintf("  %s %s (%s)\n", status, h.Args, h.Duration))
		}
	}

	b.WriteString("\n")
	b.WriteString(components.RenderFooter(
		components.FooterItem{Key: "←/Esc", Label: "Back"},
		components.FooterItem{Key: "Enter", Label: "Run"},
	))

	return b.String()
}
