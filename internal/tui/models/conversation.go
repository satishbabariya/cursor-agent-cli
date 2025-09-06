package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui/styles"
)

// ConversationModel represents the conversation view model
type ConversationModel struct {
	viewport viewport.Model
	ready    bool
}

// NewConversationModel creates a new conversation model
func NewConversationModel() ConversationModel {
	return ConversationModel{}
}

// Update updates the conversation model
func (m *ConversationModel) Update(msg tea.Msg) (ConversationModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-8)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 8
		}

	case ConversationMsg:
		content := m.renderConversation(&msg.Conversation)
		m.viewport.SetContent(content)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return *m, cmd
}

// View renders the conversation view
func (m *ConversationModel) View(width, height int, conversation *client.ConversationResponse, agent *client.Agent, errorMsg string) string {
	if agent == nil {
		return styles.ErrorStyle.Render("No agent selected")
	}

	var content strings.Builder

	// Header
	header := styles.HeaderStyle.Width(width - 4).Render(fmt.Sprintf("💬 Conversation: %s", agent.Name))
	content.WriteString(header + "\n\n")

	// Error message
	if errorMsg != "" {
		content.WriteString(styles.ErrorStyle.Render("Error: "+errorMsg) + "\n\n")
		content.WriteString(styles.HelpStyle.Render("Esc: Back | q: Quit"))
		return styles.BaseStyle.Width(width).Height(height).Render(content.String())
	}

	// Initialize viewport if not ready
	if !m.ready {
		m.viewport = viewport.New(width-4, height-8)
		m.ready = true
		if conversation != nil {
			conversationContent := m.renderConversation(conversation)
			m.viewport.SetContent(conversationContent)
		}
	} else {
		// Update viewport size
		m.viewport.Width = width - 4
		m.viewport.Height = height - 8
	}

	// Conversation content
	if conversation == nil {
		content.WriteString(styles.InfoStyle.Render("Loading conversation..."))
	} else if len(conversation.Messages) == 0 {
		content.WriteString(styles.InfoStyle.Render("No messages in this conversation."))
	} else {
		content.WriteString(m.viewport.View())
	}

	content.WriteString("\n\n")

	// Help
	helpText := "↑/↓: Scroll | f: Follow-up | Esc: Back | q: Quit"
	content.WriteString(styles.HelpStyle.Render(helpText))

	return styles.BaseStyle.Width(width).Height(height).Render(content.String())
}

// renderConversation renders the conversation messages
func (m ConversationModel) renderConversation(conversation *client.ConversationResponse) string {
	if len(conversation.Messages) == 0 {
		return styles.InfoStyle.Render("No messages in this conversation.")
	}

	var content strings.Builder

	for i, message := range conversation.Messages {
		// Message header
		emoji := m.getMessageEmoji(message.Type)
		header := fmt.Sprintf("%s Message %d", emoji, i+1)
		content.WriteString(styles.TitleStyle.Render(header) + "\n")

		// Message content
		messageLines := strings.Split(message.Text, "\n")
		for _, line := range messageLines {
			if strings.TrimSpace(line) == "" {
				content.WriteString("\n")
				continue
			}

			// Apply different styling based on message type
			var styledLine string
			switch message.Type {
			case "user_message":
				styledLine = styles.InfoStyle.Render("  " + line)
			case "agent_message":
				styledLine = styles.TableCellStyle.Render("  " + line)
			case "system_message":
				styledLine = styles.HelpStyle.Render("  " + line)
			default:
				styledLine = styles.TableCellStyle.Render("  " + line)
			}

			content.WriteString(styledLine + "\n")
		}

		// Add separator between messages
		if i < len(conversation.Messages)-1 {
			content.WriteString("\n" + strings.Repeat("─", 50) + "\n\n")
		}
	}

	return content.String()
}

// getMessageEmoji returns the appropriate emoji for a message type
func (m ConversationModel) getMessageEmoji(messageType string) string {
	switch messageType {
	case "user_message":
		return "👤"
	case "agent_message":
		return "🤖"
	case "system_message":
		return "⚙️"
	default:
		return "💬"
	}
}
