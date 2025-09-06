package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui/styles"
)

// SettingsModel represents the settings view model
type SettingsModel struct {
	selectedOption int
	options        []string
}

// NewSettingsModel creates a new settings model
func NewSettingsModel() SettingsModel {
	return SettingsModel{
		options: []string{
			"Auto-refresh agents",
			"Refresh interval",
			"Show expired agents",
			"API key info",
		},
	}
}

// Update updates the settings model
func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedOption > 0 {
				m.selectedOption--
			}
		case "down", "j":
			if m.selectedOption < len(m.options)-1 {
				m.selectedOption++
			}
		case "enter", " ":
			// Handle option selection
			switch m.selectedOption {
			case 0: // Auto-refresh toggle
				// This would be handled by the main model
			case 1: // Refresh interval
				// This would open a sub-menu or input
			case 2: // Show expired agents
				// This would toggle the setting
			case 3: // API key info
				// This would show API key information
			}
		}
	}

	return m, nil
}

// View renders the settings view
func (m SettingsModel) View(width, height int, autoRefresh bool, errorMsg string) string {
	var content strings.Builder

	// Header
	header := styles.HeaderStyle.Width(width - 4).Render("⚙️  Settings")
	content.WriteString(header + "\n\n")

	// Error message
	if errorMsg != "" {
		content.WriteString(styles.ErrorStyle.Render("Error: "+errorMsg) + "\n\n")
	}

	// Settings options
	content.WriteString(styles.TitleStyle.Render("Configuration") + "\n\n")

	for i, option := range m.options {
		var line string

		// Add current value for some options
		switch i {
		case 0: // Auto-refresh
			status := "Disabled"
			if autoRefresh {
				status = "Enabled"
			}
			line = fmt.Sprintf("%s: %s", option, status)
		case 1: // Refresh interval
			line = fmt.Sprintf("%s: 30 seconds", option)
		case 2: // Show expired agents
			line = fmt.Sprintf("%s: Dashboard setting", option)
		case 3: // API key info
			line = option
		default:
			line = option
		}

		// Style based on selection
		if i == m.selectedOption {
			line = styles.ButtonActiveStyle.Render("▶ " + line)
		} else {
			line = styles.TableCellStyle.Render("  " + line)
		}

		content.WriteString(line + "\n")
	}

	content.WriteString("\n")

	// Additional info
	content.WriteString(styles.InfoStyle.Render("Configuration is stored in ~/.cursor-cli.yaml") + "\n\n")

	// Help
	helpText := "↑/↓: Navigate | Enter/Space: Select | Esc: Back | q: Quit"
	content.WriteString(styles.HelpStyle.Render(helpText))

	return styles.BaseStyle.Width(width).Height(height).Render(content.String())
}
