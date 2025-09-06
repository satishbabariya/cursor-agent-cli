package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui/models"
)

// Run starts the TUI application
func Run(apiClient *client.Client) error {
	model := models.NewModel(apiClient)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	// Handle any final cleanup if needed
	_ = finalModel

	return nil
}
