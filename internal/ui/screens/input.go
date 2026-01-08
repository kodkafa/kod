package screens

import (
	"fmt"
	"kodkafa/internal/app/usecases"
	"kodkafa/internal/ui/components"
	"kodkafa/internal/ui/tea"

	"kodkafa/internal/ui/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea_pkg "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputModel struct {
	textInput textinput.Model
	addUC     *usecases.AddPluginUseCase
	mode      tea.InputModeType
	err       error
	loading   bool
}

func NewInputModel(addUC *usecases.AddPluginUseCase, mode tea.InputModeType) *InputModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40

	if mode == tea.InputModePath {
		ti.Placeholder = "./samples/hello-py"
	} else {
		ti.Placeholder = "plugin-name"
	}

	return &InputModel{
		textInput: ti,
		addUC:     addUC,
		mode:      mode,
	}
}

func (m *InputModel) Init() tea_pkg.Cmd {
	return textinput.Blink
}

func (m *InputModel) Update(msg tea_pkg.Msg) (tea_pkg.Model, tea_pkg.Cmd) {
	var cmd tea_pkg.Cmd

	switch msg := msg.(type) {
	case tea_pkg.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea_pkg.Msg {
				return tea.SwitchStateMsg{State: tea.StateCommandMenu}
			}
		case "enter":
			val := m.textInput.Value()
			if val != "" {
				if m.mode == tea.InputModePath {
					m.loading = true
					return m, func() tea_pkg.Msg {
						res, err := m.addUC.Execute(usecases.AddPluginInput{Source: val})
						if err != nil {
							return tea.ErrMsg{Err: err}
						}
						return tea.PluginAddedMsg{PluginName: res.Plugin.Name}
					}
				} else if m.mode == tea.InputModeName {
					return m, func() tea_pkg.Msg {
						// Pass the name. Root Model will know the context from the state machine.
						return tea.PluginSelectedMsg{PluginName: val}
					}
				}
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

func (m *InputModel) View() string {
	if m.loading {
		return "\n  Adding plugin..."
	}

	var errStr string
	if m.err != nil {
		errStr = lipgloss.NewStyle().Foreground(theme.Error).Render(fmt.Sprintf("\n  Error: %v", m.err))
	}

	prompt := "Enter Plugin Path:"
	if m.mode == tea.InputModeName {
		prompt = "Enter Plugin Name:"
	}

	return fmt.Sprintf(
		"\n %s\n\n %s\n\n %s",
		lipgloss.NewStyle().Foreground(theme.Highlight).Bold(true).Render("COMMAND INPUT"),
		prompt,
		m.textInput.View(),
	) + errStr + "\n\n " + components.RenderFooter(
		components.FooterItem{Key: "Esc", Label: "Cancel"},
		components.FooterItem{Key: "Enter", Label: "Run"},
	)
}
