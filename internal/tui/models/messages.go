package models

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
)

// AgentsMsg represents a message containing agents data
type AgentsMsg struct {
	Agents []client.Agent
}

// ConversationMsg represents a message containing conversation data
type ConversationMsg struct {
	Conversation client.ConversationResponse
}

// ErrorMsg represents an error message
type ErrorMsg struct {
	Error string
}

// TickMsg represents a tick for periodic updates
type TickMsg time.Time

// AgentSelectedMsg represents an agent selection
type AgentSelectedMsg struct {
	Index int
}

// FollowupSentMsg represents a successful followup message
type FollowupSentMsg struct {
	AgentID string
	Message string
}

// fetchAgents fetches agents from the API
func (m Model) fetchAgents() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		response, err := m.client.ListAgents(100, "")
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}
		return AgentsMsg{Agents: response.Agents}
	})
}

// fetchConversation fetches conversation for a specific agent
func (m Model) fetchConversation(agentID string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		conversation, err := m.client.GetAgentConversation(agentID)
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}
		return ConversationMsg{Conversation: *conversation}
	})
}

// sendFollowup sends a followup message to an agent
func (m Model) sendFollowup(agentID, message string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		_, err := m.client.AddFollowup(agentID, message)
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}
		return FollowupSentMsg{AgentID: agentID, Message: message}
	})
}

// tickCmd returns a command that sends a tick message every second
func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
