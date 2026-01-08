package screens

import (
	"fmt"
	"kodkafa/internal/app/dto"
	"kodkafa/internal/app/usecases"
	"kodkafa/internal/build"
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"
	"strings"

	"kodkafa/internal/ui/theme"

	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			MarginBottom(1)

	taglineStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(theme.TextSecondary)

	urlStyle = lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Underline(true)

	paginationStyle = lipgloss.NewStyle().
			Foreground(theme.TextSecondary)

	sectionHeaderStyle = lipgloss.NewStyle().
				Foreground(theme.SectionHeader)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(theme.Primary).
				Background(theme.Subtle).
				PaddingLeft(1)

	unselectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1)
)

// DashboardModel handles the main plugin list view.
type DashboardModel struct {
	listUC      *usecases.ListPluginsUseCase
	data        dto.DashboardDTO
	cursor      int
	listSource  int // 0 for Top list, 1 for Main list
	loading     bool
	err         error
	filter      string
	isSearching bool
}

// NewDashboardModel creates a new DashboardModel.
func NewDashboardModel(listUC *usecases.ListPluginsUseCase) *DashboardModel {
	return &DashboardModel{
		listUC:  listUC,
		loading: true,
	}
}

// Init fetches the initial plugin list.
func (m *DashboardModel) Init() tea_pkg.Cmd {
	return func() tea_pkg.Msg {
		data, err := m.listUC.Execute(usecases.ListPluginsInput{Page: 1, PageSize: 0})
		if err != nil {
			return tea.ErrMsg{Err: err}
		}
		return tea.PluginsLoadedMsg{Data: data}
	}
}

// Update handles input and data loading for the dashboard.
func (m *DashboardModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	switch msg := msg.(type) {
	case tea.PluginsLoadedMsg:
		m.data = msg.Data
		m.loading = false
		m.cursor = 0 // Reset cursor to avoid out-of-bounds
		return m, nil

	case tea.ErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case tea_pkg.KeyMsg:
		if m.loading {
			return m, nil
		}

		pName := m.getSelectedName()

		// Global keys if NOT searching (or handled specially)
		if !m.isSearching {
			switch msg.String() {
			case "m":
				return m, func() tea_pkg.Msg {
					return tea.SwitchStateMsg{State: tea.StateCommandMenu}
				}
			case "i":
				if pName != "" {
					return m, func() tea_pkg.Msg {
						return tea.PluginSelectedMsg{PluginName: pName, Cmd: "kod info"}
					}
				}
			case "l":
				if pName != "" {
					return m, func() tea_pkg.Msg {
						return tea.PluginSelectedMsg{PluginName: pName, Cmd: "kod load"}
					}
				}
			case "s":
				m.isSearching = true
				return m, nil
			}
		}

		// Search input handling
		if m.isSearching {
			switch msg.String() {
			case "esc":
				m.isSearching = false
				m.filter = ""
				m.updateFilter()
				return m, nil
			case "enter":
				m.isSearching = false
				return m, nil
			case "backspace":
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
					m.updateFilter()
				}
				return m, nil
			default:
				// If simple character, append
				if len(msg.String()) == 1 {
					m.filter += msg.String()
					m.updateFilter()
				}
				return m, nil
			}
		}

		// Navigation (Only when not searching, or user confirmed search)
		switch msg.String() {
		case "up":
			m.cursor--
			if m.cursor < 0 {
				if m.listSource == 1 && len(m.filteredTop()) > 0 {
					m.listSource = 0
					m.cursor = len(m.filteredTop()) - 1
				} else {
					m.cursor = 0
				}
			}
		case "down":
			m.cursor++
			if m.listSource == 0 {
				if m.cursor >= len(m.filteredTop()) {
					m.listSource = 1
					m.cursor = 0
				}
			} else {
				if m.cursor >= len(m.filteredMain()) {
					m.cursor = len(m.filteredMain()) - 1
				}
			}
		case "left":
			if m.listSource == 1 && m.data.CurrentPage > 1 {
				m.loading = true
				return m, m.fetchPage(m.data.CurrentPage - 1)
			}
		case "right":
			if m.listSource == 1 && m.data.CurrentPage < m.data.TotalPages {
				m.loading = true
				return m, m.fetchPage(m.data.CurrentPage + 1)
			}
		case "enter":
			if pName != "" {
				return m, func() tea_pkg.Msg {
					return tea.PluginSelectedMsg{PluginName: pName}
				}
			}
		}
	}
	return m, nil
}

func (m *DashboardModel) updateFilter() {
	m.cursor = 0
}

func (m *DashboardModel) filteredTop() []dto.PluginListItem {
	if m.filter == "" {
		return m.data.TopPlugins
	}
	var res []dto.PluginListItem
	for _, p := range m.data.TopPlugins {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(m.filter)) {
			res = append(res, p)
		}
	}
	return res
}

func (m *DashboardModel) filteredMain() []dto.PluginListItem {
	if m.filter == "" {
		return m.data.MainPlugins
	}
	var res []dto.PluginListItem
	for _, p := range m.data.MainPlugins {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(m.filter)) {
			res = append(res, p)
		}
	}
	return res
}

func (m *DashboardModel) getSelectedName() string {
	top := m.filteredTop()
	main := m.filteredMain()

	if m.listSource == 0 && m.cursor < len(top) {
		return top[m.cursor].Name
	} else if m.listSource == 1 && m.cursor < len(main) && m.cursor >= 0 {
		return main[m.cursor].Name
	}
	return ""
}

func (m *DashboardModel) fetchPage(page int) tea_pkg.Cmd {
	return func() tea_pkg.Msg {
		data, err := m.listUC.Execute(usecases.ListPluginsInput{Page: page, PageSize: 0})
		if err != nil {
			return tea.ErrMsg{Err: err}
		}
		return tea.PluginsLoadedMsg{Data: data}
	}
}

// View renders the dashboard.
func (m *DashboardModel) View() string {
	if m.loading {
		return "\n\n  Loading plugins..."
	}
	if m.err != nil {
		return fmt.Sprintf("\n\n  Error: %v", m.err)
	}

	var b strings.Builder

	// Header
	b.WriteString(components.RenderHeader("", build.Url))

	// Search bar
	if m.isSearching || m.filter != "" {
		b.WriteString("Search: " + m.filter)
		if m.isSearching {
			b.WriteString("█")
		}
		b.WriteString("\n\n")
	}

	line := "|───────────────────────────────────────────────────────────────────────"

	// Top List (Fixed)
	topPlugins := m.filteredTop()
	if len(topPlugins) > 0 {
		b.WriteString(sectionHeaderStyle.Render("Recents") + " ")
		b.WriteString(sectionHeaderStyle.Render(line) + "\n\n")
		for i, p := range topPlugins {
			style := unselectedItemStyle
			prefix := "  "
			if m.listSource == 0 && i == m.cursor {
				style = selectedItemStyle
				prefix = "> "
			}
			b.WriteString(style.Render(fmt.Sprintf("%s%-20s %s", prefix, p.Name, p.Description)) + "\n")
		}
	}

	// Main Inventory
	mainPlugins := m.filteredMain()
	b.WriteString("\n")
	b.WriteString(sectionHeaderStyle.Render("Plugins") + " ")
	b.WriteString(sectionHeaderStyle.Render(line) + "\n\n")
	for i, p := range mainPlugins {
		style := unselectedItemStyle
		prefix := "  "
		if m.listSource == 1 && i == m.cursor {
			style = selectedItemStyle
			prefix = "> "
		}
		b.WriteString(style.Render(fmt.Sprintf("%s%-20s %s", prefix, p.Name, p.Description)) + "\n")
	}

	b.WriteString("\n")
	pageInfo := fmt.Sprintf("Page %d of %d", m.data.CurrentPage, m.data.TotalPages)
	b.WriteString(paginationStyle.Render(pageInfo) + "\n")

	b.WriteString(components.RenderFooter(
		components.FooterItem{Key: "↑/↓", Label: ""},
		components.FooterItem{Key: "←/→", Label: ""},
		components.FooterItem{Key: "S", Label: "Search"},
		components.FooterItem{Key: "M", Label: "Menu"},
		components.FooterItem{Key: "I", Label: "Info"},
		components.FooterItem{Key: "L", Label: "Load"},
		components.FooterItem{Key: "Q", Label: "Quit"},
	))

	return b.String()
}
