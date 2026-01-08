package theme

import "github.com/charmbracelet/lipgloss"

var (
	Pink   = lipgloss.Color("#EE6983")
	Cyan   = lipgloss.Color("#51C4D3")
	Green  = lipgloss.Color("#4EDC69")
	Yellow = lipgloss.Color("#FFD93D")
	Subtle = lipgloss.Color("#222222")

	// Main Palette
	Primary   = Pink
	Secondary = Cyan
	Accent    = Yellow
	Muted     = lipgloss.Color("#666666")

	// Text Colors
	TextPrimary   = lipgloss.Color("#E2E2E2")
	TextSecondary = lipgloss.Color("#AAAAAA")
	TextDark      = lipgloss.Color("#171717") // For use on bright backgrounds

	// UI Elements
	SectionHeader = lipgloss.Color("#535353")
	HeaderTitle   = lipgloss.Color("#737373")
	FooterBg      = lipgloss.Color("#333333")

	// Status Colors
	Success = lipgloss.Color("#4EDC69")
	Error   = lipgloss.Color("#FF6B6B")

	// Special
	Highlight    = lipgloss.Color("#00D7FF")
	Logo         = lipgloss.Color("#4ECDC4")
	LogoGradient = []string{"#00d492"}
)
