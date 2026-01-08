package screens

import (
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"
	"strings"

	"kodkafa/internal/ui/theme"

	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	menuStyle = lipgloss.NewStyle().
			Width(80)

	menuItemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	menuSelectedStyle = lipgloss.NewStyle().
				Foreground(theme.TextDark).
				Background(theme.Primary).
				PaddingLeft(2)
)

type MenuItem struct {
	Label string
	Usage string
	Cmd   string
}

type MenuModel struct {
	items  []MenuItem
	cursor int
}

func NewMenuModel() *MenuModel {
	return &MenuModel{
		items: []MenuItem{
			{
				Label: "kodkafa add <path|repo_url> (kod a <path|repo_url>)",
				Usage: ": Add new plugin",
				Cmd:   "kod add",
			},
			{
				Label: "kodkafa run <name> (kod r <name>)",
				Usage: ": Run plugin with memory",
				Cmd:   "kod run",
			},
			{
				Label: "kodkafa info <name> (kod i <name>)",
				Usage: ": View plugin metadata & history",
				Cmd:   "kod info",
			},
			{
				Label: "kodkafa del <name> (kod d <name>)",
				Usage: ": Delete plugin safely",
				Cmd:   "kod del",
			},
			{
				Label: "kodkafa load <name> (kod l <name>)",
				Usage: ": Install/Refresh dependencies",
				Cmd:   "kod load",
			},
			{
				Label: "kodkafa log <name>",
				Usage: ": View execution logs",
				Cmd:   "kod log",
			},
			{
				Label: "kodkafa init",
				Usage: ": Initialize system layout",
				Cmd:   "kod init",
			},
		},
	}
}

func (m *MenuModel) Init() tea_pkg.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	switch msg := msg.(type) {
	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.items) - 1
			}
		case "down":
			m.cursor++
			if m.cursor >= len(m.items) {
				m.cursor = 0
			}
		case "esc", "left":
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateNormal}
			}
		case "enter", "right":
			item := m.items[m.cursor]
			switch item.Cmd {
			case "kod add":
				return m, func() tea_pkg.Msg {
					return tea.SwitchStateMsg{State: tea.StateInput, Mode: tea.InputModePath, Cmd: "kod add"}
				}
			case "kod del":
				return m, func() tea_pkg.Msg {
					return tea.SwitchStateMsg{State: tea.StateInput, Mode: tea.InputModeName, Cmd: "kod del"}
				}
			case "kod load":
				return m, func() tea_pkg.Msg {
					return tea.SwitchStateMsg{State: tea.StateInput, Mode: tea.InputModeName, Cmd: "kod load"}
				}
			case "kod run":
				return m, func() tea_pkg.Msg {
					return tea.SwitchStateMsg{State: tea.StateInput, Mode: tea.InputModeName, Cmd: "kod run"}
				}
			case "kod info":
				return m, func() tea_pkg.Msg {
					return tea.SwitchStateMsg{State: tea.StateInput, Mode: tea.InputModeName, Cmd: "kod info"}
				}
			case "kod log":
				return m, func() tea_pkg.Msg {
					return tea.SwitchStateMsg{State: tea.StateInput, Mode: tea.InputModeName, Cmd: "kod log"}
				}
			case "kod init":
				return m, func() tea_pkg.Msg {
					return tea.InitLayoutMsg{}
				}
			}
		}
	}
	return m, nil
}

func (m *MenuModel) View() string {
	var b strings.Builder

	// Header
	b.WriteString(components.RenderHeader("A Persistent CLI with Memory", "CLI COMMANDS"))

	for i, item := range m.items {
		style := menuItemStyle
		if i == m.cursor {
			style = menuSelectedStyle
		}
		b.WriteString(style.Render(item.Label+item.Usage) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(components.RenderFooter(
		components.FooterItem{Key: "↑/↓", Label: ""},
		components.FooterItem{Key: "←/Esc", Label: "Back"},
		components.FooterItem{Key: "Enter/→", Label: "Run"},
	))

	return b.String()
}
