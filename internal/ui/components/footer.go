package components

import (
	"strings"

	"kodkafa/internal/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

var (
	shortcutStyle = lipgloss.NewStyle().
			Background(theme.FooterBg).
			Foreground(theme.TextPrimary).
			PaddingLeft(1).
			PaddingRight(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(theme.TextSecondary).
			PaddingLeft(1).
			PaddingRight(2)
)

// FooterItem represents a single action in the footer
type FooterItem struct {
	Key   string
	Label string
}

// RenderFooter renders a horizontal list of shortcut keys and labels.
func RenderFooter(items ...FooterItem) string {
	var renderedItems []string

	for _, item := range items {
		runes := []rune(item.Label)
		// If key matches first letter of label, render compactly (e.g. "S Search" -> "Search" with built-in styling)
		if len(runes) > 0 && strings.EqualFold(item.Key, string(runes[0])) {
			renderedItems = append(renderedItems, shortcutStyle.Copy().PaddingRight(1).Render(string(runes[0])))
			renderedItems = append(renderedItems, labelStyle.Copy().PaddingLeft(0).Render(string(runes[1:])))
		} else {
			renderedItems = append(renderedItems, shortcutStyle.Render(item.Key))
			renderedItems = append(renderedItems, labelStyle.Render(item.Label))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedItems...)
}
