package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Success   = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Error     = lipgloss.Color("#EF4444") // Red
	Muted     = lipgloss.Color("#6B7280") // Gray

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FBBF24")).
			Background(Primary).
			Bold(true).
			Padding(0, 1)

	// Status styles
	StatusRunning = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	StatusCompleted = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	StatusFailed = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	StatusExpired = lipgloss.NewStyle().
			Foreground(Muted).
			Bold(true)

	StatusCancelled = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FBBF24")).
				Background(Primary).
				Bold(true).
				Padding(0, 1).
				Margin(0)

	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0)

	TableSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(Secondary).
				Bold(true).
				Padding(0, 1).
				Margin(0)

	// Border styles
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Margin(1, 0)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Margin(0, 0, 1, 0)

	// Error message styles
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true).
			Margin(1, 0)

	// Success message styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true).
			Margin(1, 0)

	// Info message styles
	InfoStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true).
			Margin(1, 0)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(Secondary).
				Padding(0, 1)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Background(Primary).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2).
			Margin(0, 1)

	ButtonActiveStyle = lipgloss.NewStyle().
				Background(Secondary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)
)

// GetStatusStyle returns the appropriate style for an agent status
func GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "RUNNING":
		return StatusRunning
	case "COMPLETED":
		return StatusCompleted
	case "FAILED":
		return StatusFailed
	case "EXPIRED":
		return StatusExpired
	case "CANCELLED":
		return StatusCancelled
	default:
		return TableCellStyle
	}
}

// GetStatusEmoji returns the appropriate emoji for an agent status
func GetStatusEmoji(status string) string {
	switch status {
	case "RUNNING":
		return "🏃"
	case "COMPLETED":
		return "✅"
	case "FAILED":
		return "❌"
	case "EXPIRED":
		return "⏰"
	case "CANCELLED":
		return "🚫"
	default:
		return "❓"
	}
}
