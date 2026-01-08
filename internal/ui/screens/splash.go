package screens

import (
	"kodkafa/internal/build"
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"
	"kodkafa/internal/ui/theme"
	"time"

	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SplashModel handles the animated splash screen.
type SplashModel struct {
	lines       []string
	currentLine int
	done        bool
}

// NewSplashModel creates a new SplashModel.
func NewSplashModel() *SplashModel {
	return &SplashModel{
		lines: components.RenderLogoLines(),
	}
}

// Init starts the animation tick.
func (m *SplashModel) Init() tea_pkg.Cmd {
	return m.tick()
}

func (m *SplashModel) tick() tea_pkg.Cmd {
	return tea_pkg.Tick(time.Millisecond*100, func(t time.Time) tea_pkg.Msg {
		return tea.TickMsg{}
	})
}

// Update handles animation ticks and transitions.
func (m *SplashModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	switch msg.(type) {
	case tea.TickMsg:
		if m.currentLine < len(m.lines) {
			m.currentLine++
			return m, m.tick()
		}
		m.done = true
		// Transition to dashboard after a short delay
		return m, tea_pkg.Tick(time.Second, func(t time.Time) tea_pkg.Msg {
			return tea.SwitchScreenMsg{ScreenName: "dashboard"}
		})
	}
	return m, nil
}

// View renders the partially drawn logo.
func (m *SplashModel) View() string {
	var out string
	for i := range m.lines {
		if i < m.currentLine {
			out += m.lines[i] + "\n"
		} else {
			out += "\n"
		}
	}
	// Center vertically or just add some padding

	out = lipgloss.NewStyle().MarginTop(2).Render(out)
	out += "\n"
	out += lipgloss.NewStyle().Foreground(theme.HeaderTitle).Bold(true).Render(build.Tagline)
	out += "\n"
	out += lipgloss.NewStyle().Foreground(theme.HeaderTitle).Render(build.Repo)
	out += "\n"
	out += lipgloss.NewStyle().Foreground(theme.HeaderTitle).Render(build.Url)
	out += "\n\n\n"
	return out
}
