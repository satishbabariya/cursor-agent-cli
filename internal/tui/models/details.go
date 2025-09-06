package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui/styles"
)

// DetailsModel represents the agent details view model
type DetailsModel struct {
	scrollOffset int
}

// NewDetailsModel creates a new details model
func NewDetailsModel() DetailsModel {
	return DetailsModel{}
}

// Update updates the details model
func (m DetailsModel) Update(msg tea.Msg) (DetailsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case "down", "j":
			m.scrollOffset++
		case "home":
			m.scrollOffset = 0
		case "end":
			m.scrollOffset = 100 // Will be clamped in View
		}
	}

	return m, nil
}

// View renders the details view
func (m DetailsModel) View(width, height int, agent *client.Agent, errorMsg string) string {
	if agent == nil {
		return styles.ErrorStyle.Render("No agent selected")
	}

	var content strings.Builder

	// Header
	header := styles.HeaderStyle.Width(width - 4).Render(fmt.Sprintf("🤖 Agent Details: %s", agent.Name))
	content.WriteString(header + "\n\n")

	// Error message
	if errorMsg != "" {
		content.WriteString(styles.ErrorStyle.Render("Error: "+errorMsg) + "\n\n")
	}

	// Agent information sections
	sections := []string{
		m.renderBasicInfo(agent),
		m.renderSourceInfo(agent),
		m.renderTargetInfo(agent),
		m.renderSummaryInfo(agent),
	}

	fullContent := strings.Join(sections, "\n\n")
	lines := strings.Split(fullContent, "\n")

	// Handle scrolling
	availableHeight := height - 6 // Account for header, help, padding
	if len(lines) > availableHeight {
		start := m.scrollOffset
		end := start + availableHeight
		if end > len(lines) {
			end = len(lines)
			start = end - availableHeight
			if start < 0 {
				start = 0
			}
		}
		lines = lines[start:end]
	}

	content.WriteString(strings.Join(lines, "\n") + "\n\n")

	// Help
	helpText := "↑/↓: Scroll | c: Conversation | f: Follow-up | Esc: Back | q: Quit"
	content.WriteString(styles.HelpStyle.Render(helpText))

	return styles.BaseStyle.Width(width).Height(height).Render(content.String())
}

// renderBasicInfo renders the basic agent information
func (m DetailsModel) renderBasicInfo(agent *client.Agent) string {
	var content strings.Builder

	content.WriteString(styles.TitleStyle.Render("📋 Basic Information") + "\n")

	// Create info table
	info := [][]string{
		{"ID:", agent.ID},
		{"Name:", agent.Name},
		{"Status:", fmt.Sprintf("%s %s", styles.GetStatusEmoji(agent.Status), agent.Status)},
		{"Created:", agent.CreatedAt.Format("2006-01-02 15:04:05 MST")},
	}

	for _, row := range info {
		label := styles.TableCellStyle.Bold(true).Render(row[0])
		value := styles.TableCellStyle.Render(row[1])
		content.WriteString(fmt.Sprintf("  %s %s\n", label, value))
	}

	return content.String()
}

// renderSourceInfo renders the source repository information
func (m DetailsModel) renderSourceInfo(agent *client.Agent) string {
	var content strings.Builder

	content.WriteString(styles.TitleStyle.Render("📂 Source Information") + "\n")

	info := [][]string{
		{"Repository:", agent.Source.Repository},
		{"Reference:", agent.Source.Ref},
	}

	for _, row := range info {
		label := styles.TableCellStyle.Bold(true).Render(row[0])
		value := styles.TableCellStyle.Render(row[1])
		content.WriteString(fmt.Sprintf("  %s %s\n", label, value))
	}

	return content.String()
}

// renderTargetInfo renders the target branch and PR information
func (m DetailsModel) renderTargetInfo(agent *client.Agent) string {
	var content strings.Builder

	content.WriteString(styles.TitleStyle.Render("🎯 Target Information") + "\n")

	info := [][]string{
		{"Branch:", agent.Target.BranchName},
		{"Agent URL:", agent.Target.URL},
		{"Auto Create PR:", fmt.Sprintf("%t", agent.Target.AutoCreatePr)},
	}

	if agent.Target.PrURL != "" {
		info = append(info, []string{"Pull Request:", agent.Target.PrURL})
	}

	for _, row := range info {
		label := styles.TableCellStyle.Bold(true).Render(row[0])
		value := styles.TableCellStyle.Render(row[1])
		content.WriteString(fmt.Sprintf("  %s %s\n", label, value))
	}

	return content.String()
}

// renderSummaryInfo renders the agent summary
func (m DetailsModel) renderSummaryInfo(agent *client.Agent) string {
	if agent.Summary == "" {
		return ""
	}

	var content strings.Builder

	content.WriteString(styles.TitleStyle.Render("📄 Summary") + "\n")

	// Word wrap the summary
	summary := m.wordWrap(agent.Summary, 70)
	for _, line := range strings.Split(summary, "\n") {
		content.WriteString(fmt.Sprintf("  %s\n", styles.TableCellStyle.Render(line)))
	}

	return content.String()
}

// wordWrap wraps text to the specified width
func (m DetailsModel) wordWrap(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+len(word)+1 <= width {
			currentLine.WriteString(" " + word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}
