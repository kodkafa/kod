package components

import (
	"kodkafa/internal/build"
	"kodkafa/internal/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

var (
	taglineStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(theme.TextSecondary)

	headerTitleStyle = lipgloss.NewStyle().
				Foreground(theme.HeaderTitle)

	separatorStyle = lipgloss.NewStyle().
			Foreground(theme.Accent).
			BorderBottom(true).
			MarginBottom(1)
)

// RenderHeader renders a standard header with the small logo, tagline, and screen title.
// Format:
// SMALL LOGO - {tagline}
// {title}
// ---
func RenderHeader(tagline, title string) string {
	if tagline == "" {
		tagline = build.Tagline
	}
	logo := GetSmallLogo()

	headerTop := logo
	if tagline != "" {
		headerTop += " " + taglineStyle.Render("- "+tagline)
	}

	out := headerTop + " " + build.Version + "\n"
	if title != "" {
		out += headerTitleStyle.Render(title) + "\n"
	}

	// Solid line 80 width
	line := "────────────────────────────────────────────────────────────────────────────────"
	out += separatorStyle.Render(line) + "\n"

	return out
}
