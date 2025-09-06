package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui/styles"
)

// HelpModel represents the help view model
type HelpModel struct{}

// NewHelpModel creates a new help model
func NewHelpModel() HelpModel {
	return HelpModel{}
}

// Update updates the help model
func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	return m, nil
}

// View renders the help view
func (m HelpModel) View(width, height int, keyMap KeyMap) string {
	var content strings.Builder

	// Header
	header := styles.HeaderStyle.Width(width - 4).Render("❓ Help & Keyboard Shortcuts")
	content.WriteString(header + "\n\n")

	// Navigation section
	content.WriteString(styles.TitleStyle.Render("🧭 Navigation") + "\n")
	navHelp := [][]string{
		{"↑/k, ↓/j", "Navigate up/down"},
		{"←/h, →/l", "Navigate left/right"},
		{"Enter", "Select/confirm"},
		{"Esc", "Go back"},
		{"Tab", "Next view"},
		{"q, Ctrl+C", "Quit application"},
	}
	content.WriteString(m.renderHelpTable(navHelp) + "\n\n")

	// Dashboard section
	content.WriteString(styles.TitleStyle.Render("📋 Dashboard") + "\n")
	dashboardHelp := [][]string{
		{"r, F5", "Refresh agents"},
		{"t", "Toggle show all/active agents"},
		{"d", "View agent details"},
		{"c", "View conversation"},
		{"f", "Send follow-up (running agents)"},
		{"s", "Open settings"},
		{"?", "Show this help"},
	}
	content.WriteString(m.renderHelpTable(dashboardHelp) + "\n\n")

	// Views section
	content.WriteString(styles.TitleStyle.Render("👁️  Views") + "\n")
	viewsHelp := [][]string{
		{"Details View", "Scroll with ↑/↓, view agent information"},
		{"Conversation", "Scroll through message history"},
		{"Follow-up", "Ctrl+T to toggle input mode"},
		{"Settings", "Configure application preferences"},
	}
	content.WriteString(m.renderHelpTable(viewsHelp) + "\n\n")

	// Tips section
	content.WriteString(styles.TitleStyle.Render("💡 Tips") + "\n")
	tips := []string{
		"• Agents auto-refresh every 30 seconds",
		"• Use 't' in dashboard to filter expired agents",
		"• Follow-up messages can only be sent to running agents",
		"• Press Ctrl+T in follow-up view to switch between short and long message input",
		"• Configuration is saved in ~/.cursor-cli.yaml",
	}
	for _, tip := range tips {
		content.WriteString(styles.TableCellStyle.Render(tip) + "\n")
	}

	content.WriteString("\n")

	// Footer
	content.WriteString(styles.HelpStyle.Render("Press Esc to return to the dashboard"))

	return styles.BaseStyle.Width(width).Height(height).Render(content.String())
}

// renderHelpTable renders a help table with key bindings
func (m HelpModel) renderHelpTable(items [][]string) string {
	var content strings.Builder

	for _, item := range items {
		key := styles.ButtonStyle.Render(item[0])
		desc := styles.TableCellStyle.Render(item[1])
		content.WriteString("  " + key + " " + desc + "\n")
	}

	return content.String()
}
