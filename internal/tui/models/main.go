package models

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
)

// View represents the different views in the TUI
type View int

const (
	DashboardView View = iota
	AgentDetailsView
	ConversationView
	FollowupView
	SettingsView
	HelpView
)

// Model represents the main TUI model
type Model struct {
	// Current state
	currentView View
	width       int
	height      int

	// API client
	client *client.Client

	// Data
	agents        []client.Agent
	selectedAgent *client.Agent
	conversation  *client.ConversationResponse

	// Sub-models
	dashboard         DashboardModel
	details           DetailsModel
	conversationModel ConversationModel
	followup          FollowupModel
	settings          SettingsModel
	help              HelpModel

	// State
	loading     bool
	error       string
	lastRefresh time.Time
	autoRefresh bool

	// Key bindings
	keyMap KeyMap
}

// KeyMap defines the key bindings for the TUI
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Enter        key.Binding
	Back         key.Binding
	Refresh      key.Binding
	Help         key.Binding
	Quit         key.Binding
	Tab          key.Binding
	Details      key.Binding
	Conversation key.Binding
	Followup     key.Binding
	Settings     key.Binding
	Toggle       key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "f5"),
			key.WithHelp("r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?", "f1"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next view"),
		),
		Details: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "details"),
		),
		Conversation: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "conversation"),
		),
		Followup: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "followup"),
		),
		Settings: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "settings"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle"),
		),
	}
}

// NewModel creates a new TUI model
func NewModel(apiClient *client.Client) Model {
	m := Model{
		currentView: DashboardView,
		client:      apiClient,
		keyMap:      DefaultKeyMap(),
		autoRefresh: true,
	}

	// Initialize sub-models
	m.dashboard = NewDashboardModel()
	m.details = NewDetailsModel()
	m.conversationModel = NewConversationModel()
	m.followup = NewFollowupModel()
	m.settings = NewSettingsModel()
	m.help = NewHelpModel()

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchAgents(),
		m.tickCmd(),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.Help):
			m.currentView = HelpView

		case key.Matches(msg, m.keyMap.Back):
			if m.currentView != DashboardView {
				m.currentView = DashboardView
			}

		case key.Matches(msg, m.keyMap.Refresh):
			cmd = m.fetchAgents()
			cmds = append(cmds, cmd)

		case key.Matches(msg, m.keyMap.Enter):
			if m.currentView == DashboardView && m.selectedAgent != nil {
				m.currentView = AgentDetailsView
			}

		case key.Matches(msg, m.keyMap.Details):
			if m.selectedAgent != nil {
				m.currentView = AgentDetailsView
			}

		case key.Matches(msg, m.keyMap.Conversation):
			if m.selectedAgent != nil {
				m.currentView = ConversationView
				cmd = m.fetchConversation(m.selectedAgent.ID)
				cmds = append(cmds, cmd)
			}

		case key.Matches(msg, m.keyMap.Followup):
			if m.selectedAgent != nil && m.selectedAgent.Status == "RUNNING" {
				m.currentView = FollowupView
			}

		case key.Matches(msg, m.keyMap.Settings):
			m.currentView = SettingsView
		}

	case AgentsMsg:
		m.agents = msg.Agents
		m.loading = false
		m.error = ""
		m.lastRefresh = time.Now()

		// Update dashboard
		m.dashboard, cmd = m.dashboard.Update(msg)
		cmds = append(cmds, cmd)

	case ConversationMsg:
		m.conversation = &msg.Conversation
		m.conversationModel, cmd = m.conversationModel.Update(msg)
		cmds = append(cmds, cmd)

	case ErrorMsg:
		m.error = msg.Error
		m.loading = false

	case TickMsg:
		if m.autoRefresh && time.Since(m.lastRefresh) > 30*time.Second {
			cmd = m.fetchAgents()
			cmds = append(cmds, cmd)
		}
		cmd = m.tickCmd()
		cmds = append(cmds, cmd)

	case AgentSelectedMsg:
		m.selectedAgent = &msg.Agent
		// Automatically switch to details view when agent is selected
		m.currentView = AgentDetailsView
	}

	// Update current view
	switch m.currentView {
	case DashboardView:
		m.dashboard, cmd = m.dashboard.Update(msg)
	case AgentDetailsView:
		m.details, cmd = m.details.Update(msg)
	case ConversationView:
		m.conversationModel, cmd = m.conversationModel.Update(msg)
	case FollowupView:
		m.followup, cmd = m.followup.Update(msg)
	case SettingsView:
		m.settings, cmd = m.settings.Update(msg)
	case HelpView:
		m.help, cmd = m.help.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the current view
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	switch m.currentView {
	case DashboardView:
		return m.dashboard.View(m.width, m.height, m.agents, m.selectedAgent, m.error)
	case AgentDetailsView:
		return m.details.View(m.width, m.height, m.selectedAgent, m.error)
	case ConversationView:
		return m.conversationModel.View(m.width, m.height, m.conversation, m.selectedAgent, m.error)
	case FollowupView:
		return m.followup.View(m.width, m.height, m.selectedAgent, m.error)
	case SettingsView:
		return m.settings.View(m.width, m.height, m.autoRefresh, m.error)
	case HelpView:
		return m.help.View(m.width, m.height, m.keyMap)
	default:
		return "Unknown view"
	}
}
