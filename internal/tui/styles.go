package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary   = lipgloss.Color("#7C3AED") // purple
	colorSecondary = lipgloss.Color("#10B981") // green
	colorWarning   = lipgloss.Color("#F59E0B") // yellow
	colorDanger    = lipgloss.Color("#EF4444") // red
	colorMuted     = lipgloss.Color("#6B7280") // gray
	colorText      = lipgloss.Color("#F9FAFB") // white
	colorBg        = lipgloss.Color("#1F2937") // dark bg
	colorStar      = lipgloss.Color("#FBBF24") // gold

	// Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginBottom(1)

	// Box styles
	goalBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSecondary).
			Padding(0, 1).
			MarginBottom(1)

	editorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	// Labels
	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSecondary)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorText).
			MarginTop(1)

	modeNormalStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(colorSecondary).
			Padding(0, 1)

	modeInsertStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(colorPrimary).
			Padding(0, 1)

	modeVisualStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(colorWarning).
			Padding(0, 1)

	// Stars
	starStyle = lipgloss.NewStyle().
			Foreground(colorStar)

	noStarStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Track header in level list
	trackHeaderStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Bold(true)

	// Menu items
	selectedStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	unselectedStyle = lipgloss.NewStyle().
			Foreground(colorText)

	lockedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Help text
	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)

	// Danger/warning text
	dangerStyle = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true).
			MarginTop(1)

	// Success/clear message
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSecondary).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorSecondary).
			Padding(1, 2).
			Align(lipgloss.Center)

	// Hint
	hintStyle = lipgloss.NewStyle().
			Foreground(colorWarning).
			Italic(true).
			MarginTop(1)

	// Solution
	solutionStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			MarginTop(1)

	// Solution explanation
	explanationStyle = lipgloss.NewStyle().
				Foreground(colorText).
				Italic(true)

	// Cursor character highlight
	cursorStyle = lipgloss.NewStyle().
			Reverse(true)
)

// FormatStars returns a star display string.
func FormatStars(stars int) string {
	s := ""
	for i := 0; i < 3; i++ {
		if i < stars {
			s += starStyle.Render("*")
		} else {
			s += noStarStyle.Render("*")
		}
	}
	return s
}

// ModeStyle returns the appropriate style for a vim mode.
func ModeStyle(mode string) lipgloss.Style {
	switch mode {
	case "NORMAL":
		return modeNormalStyle
	case "INSERT":
		return modeInsertStyle
	case "VISUAL", "V-LINE", "V-BLOCK":
		return modeVisualStyle
	default:
		return modeNormalStyle
	}
}
