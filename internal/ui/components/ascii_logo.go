package components

import (
	"kodkafa/internal/build"
	"kodkafa/internal/ui/theme"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	LogoStyle = lipgloss.NewStyle().
			Foreground(theme.Logo).
			Bold(true)

	LogoSmallStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Logo)
)

// RenderLogo renders the full-size logo with optional color overrides
func RenderLogo(color string) string {
	style := LogoStyle
	if color != "" {
		style = style.Foreground(lipgloss.Color(color))
	}
	return style.Render(build.LogoAscii)
}

// RenderLogoLines returns the logo split into individual lines with a horizontal gradient
func RenderLogoLines() []string {
	lines := strings.Split(build.LogoAscii, "\n")
	rendered := make([]string, len(lines))

	// Colors for gradient: teal to blue
	colors := theme.LogoGradient
	for i, line := range lines {
		if i < len(colors) {
			rendered[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i])).Render(line)
		} else {
			rendered[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(colors[0])).Render(line)
		}
	}
	return rendered
}

func GetSmallLogo() string {
	return LogoSmallStyle.Render(build.Name)
}
