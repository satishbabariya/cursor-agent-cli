package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui/styles"
)

// DashboardModel represents the dashboard view model
type DashboardModel struct {
	table          table.Model
	showAll        bool
	selectedRow    int
	filteredAgents []client.Agent
	allAgents      []client.Agent
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel() DashboardModel {
	columns := []table.Column{
		{Title: "ID", Width: 12},
		{Title: "Name", Width: 25},
		{Title: "Status", Width: 12},
		{Title: "Repository", Width: 30},
		{Title: "Created", Width: 16},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Primary).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(styles.Secondary).
		Bold(false)
	t.SetStyles(s)

	return DashboardModel{
		table:   t,
		showAll: false,
	}
}

// Update updates the dashboard model
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("t"))):
			m.showAll = !m.showAll
			// Re-filter the table with current agents
			m.updateTable(m.allAgents)

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			selectedRow := m.table.Cursor()
			if selectedRow < len(m.filteredAgents) {
				// Find the index in the original agents slice
				selectedAgent := m.filteredAgents[selectedRow]
				return m, func() tea.Msg {
					return AgentSelectedMsg{Agent: selectedAgent}
				}
			}
		}

	case AgentsMsg:
		m.updateTable(msg.Agents)
	}

	m.table, cmd = m.table.Update(msg)
	m.selectedRow = m.table.Cursor()

	return m, cmd
}

// View renders the dashboard view
func (m DashboardModel) View(width, height int, agents []client.Agent, selectedAgent *client.Agent, errorMsg string) string {
	var content strings.Builder

	// Header
	header := styles.HeaderStyle.Width(width - 4).Render("🚀 Cursor Background Agents Dashboard")
	content.WriteString(header + "\n\n")

	// Filter info
	filterText := "● Active Agents"
	if m.showAll {
		filterText = "● All Agents"
	}

	activeCount := 0
	for _, agent := range agents {
		if agent.Status != "EXPIRED" {
			activeCount++
		}
	}

	statusLine := fmt.Sprintf("%s (%d) | Press 't' to toggle | Press '?' for help",
		filterText,
		func() int {
			if m.showAll {
				return len(agents)
			}
			return activeCount
		}())

	content.WriteString(styles.InfoStyle.Render(statusLine) + "\n")

	// Error message
	if errorMsg != "" {
		content.WriteString(styles.ErrorStyle.Render("Error: "+errorMsg) + "\n")
	}

	// Table
	m.table.SetWidth(width - 4)
	m.table.SetHeight(height - 10)
	content.WriteString(m.table.View() + "\n")

	// Help
	helpText := "↑/↓: Navigate | Enter: View Details | d: Details | c: Conversation | f: Follow-up | r: Refresh | q: Quit"
	content.WriteString(styles.HelpStyle.Render(helpText))

	return styles.BaseStyle.Width(width).Height(height).Render(content.String())
}

// updateTable updates the table with new agent data
func (m *DashboardModel) updateTable(agents []client.Agent) {
	m.allAgents = agents // Store all agents

	var rows []table.Row
	var filteredAgents []client.Agent

	for _, agent := range agents {
		// Filter expired agents if not showing all
		if !m.showAll && agent.Status == "EXPIRED" {
			continue
		}

		filteredAgents = append(filteredAgents, agent)

		statusText := fmt.Sprintf("%s %s",
			styles.GetStatusEmoji(agent.Status),
			agent.Status)

		// Extract repository name from URL
		repoName := agent.Source.Repository
		if strings.Contains(repoName, "/") {
			parts := strings.Split(repoName, "/")
			if len(parts) >= 2 {
				repoName = strings.Join(parts[len(parts)-2:], "/")
			}
		}

		row := table.Row{
			agent.ID,
			truncateString(agent.Name, 25),
			statusText,
			truncateString(repoName, 30),
			agent.CreatedAt.Format("2006-01-02 15:04"),
		}
		rows = append(rows, row)
	}

	m.filteredAgents = filteredAgents
	m.table.SetRows(rows)
}

// truncateString truncates a string to the specified length
func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
