package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui/styles"
)

// FollowupModel represents the followup message view model
type FollowupModel struct {
	textarea    textarea.Model
	textinput   textinput.Model
	useTextarea bool
	sending     bool
	sent        bool
	error       string
}

// NewFollowupModel creates a new followup model
func NewFollowupModel() FollowupModel {
	// Text input for short messages
	ti := textinput.New()
	ti.Placeholder = "Enter your follow-up message..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	// Text area for longer messages
	ta := textarea.New()
	ta.Placeholder = "Enter your follow-up message (press Ctrl+Enter to send)..."
	ta.SetWidth(70)
	ta.SetHeight(8)

	return FollowupModel{
		textarea:    ta,
		textinput:   ti,
		useTextarea: false,
	}
}

// Update updates the followup model
func (m FollowupModel) Update(msg tea.Msg) (FollowupModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+t":
			// Toggle between text input and textarea
			m.useTextarea = !m.useTextarea
			if m.useTextarea {
				m.textarea.Focus()
				m.textinput.Blur()
			} else {
				m.textinput.Focus()
				m.textarea.Blur()
			}

		case "enter":
			if !m.useTextarea && m.textinput.Value() != "" && !m.sending {
				return m, m.sendMessage(m.textinput.Value())
			}

		case "ctrl+enter":
			if m.useTextarea && m.textarea.Value() != "" && !m.sending {
				return m, m.sendMessage(m.textarea.Value())
			}

		case "esc":
			if !m.sending {
				m.reset()
			}
		}

	case FollowupSentMsg:
		m.sending = false
		m.sent = true
		m.error = ""

	case ErrorMsg:
		m.sending = false
		m.error = msg.Error
	}

	// Update the appropriate input
	if m.useTextarea {
		m.textarea, cmd = m.textarea.Update(msg)
	} else {
		m.textinput, cmd = m.textinput.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the followup view
func (m FollowupModel) View(width, height int, agent *client.Agent, errorMsg string) string {
	if agent == nil {
		return styles.ErrorStyle.Render("No agent selected")
	}

	if agent.Status != "RUNNING" {
		return styles.ErrorStyle.Render("Can only send follow-up messages to running agents")
	}

	var content strings.Builder

	// Header
	header := styles.HeaderStyle.Width(width - 4).Render(fmt.Sprintf("📤 Send Follow-up: %s", agent.Name))
	content.WriteString(header + "\n\n")

	// Success message
	if m.sent {
		content.WriteString(styles.SuccessStyle.Render("✅ Follow-up message sent successfully!") + "\n\n")
		content.WriteString(styles.InfoStyle.Render("The agent will process your instruction and continue working.") + "\n\n")
	}

	// Error message
	if m.error != "" {
		content.WriteString(styles.ErrorStyle.Render("Error: "+m.error) + "\n\n")
	}
	if errorMsg != "" {
		content.WriteString(styles.ErrorStyle.Render("Error: "+errorMsg) + "\n\n")
	}

	// Instructions
	if !m.sent {
		content.WriteString(styles.InfoStyle.Render("Enter additional instructions for the agent:") + "\n\n")

		// Input mode toggle
		inputMode := "Short message"
		if m.useTextarea {
			inputMode = "Long message"
		}
		content.WriteString(styles.HelpStyle.Render(fmt.Sprintf("Mode: %s (Ctrl+T to toggle)", inputMode)) + "\n\n")

		// Input field
		if m.useTextarea {
			content.WriteString(m.textarea.View() + "\n")
		} else {
			content.WriteString(m.textinput.View() + "\n\n")
		}

		// Sending status
		if m.sending {
			content.WriteString(styles.InfoStyle.Render("📤 Sending follow-up message...") + "\n")
		}
	}

	// Help
	var helpText string
	if m.sent {
		helpText = "Esc: Back | q: Quit"
	} else if m.useTextarea {
		helpText = "Ctrl+Enter: Send | Ctrl+T: Toggle input mode | Esc: Back | q: Quit"
	} else {
		helpText = "Enter: Send | Ctrl+T: Toggle input mode | Esc: Back | q: Quit"
	}
	content.WriteString("\n" + styles.HelpStyle.Render(helpText))

	return styles.BaseStyle.Width(width).Height(height).Render(content.String())
}

// sendMessage creates a command to send a followup message
func (m FollowupModel) sendMessage(message string) tea.Cmd {
	m.sending = true
	return func() tea.Msg {
		// This will be handled by the main model
		return FollowupSentMsg{Message: message}
	}
}

// reset resets the followup model state
func (m *FollowupModel) reset() {
	m.textinput.SetValue("")
	m.textarea.SetValue("")
	m.sending = false
	m.sent = false
	m.error = ""
}
