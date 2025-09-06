package cmd

import (
	"fmt"
	"os"

	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/config"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/tui"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive Text User Interface",
	Long: `Launch the interactive Text User Interface (TUI) for managing Cursor Background Agents.

The TUI provides a rich, interactive experience with:
- Real-time agent dashboard with status updates
- Detailed agent information and progress monitoring  
- Conversation history viewer with syntax highlighting
- Interactive follow-up message composition
- Keyboard shortcuts for efficient navigation
- Auto-refresh capabilities

This is perfect for monitoring multiple agents and their progress in real-time.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := config.GetAPIKey()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			fmt.Println("💡 Run 'cursor-cli init' to set up your API key first.")
			os.Exit(1)
		}

		client := client.NewClient(apiKey)

		fmt.Println("🚀 Starting Cursor Background Agents TUI...")
		fmt.Println("💡 Press '?' for help, 'q' to quit")

		if err := tui.Run(client); err != nil {
			fmt.Printf("❌ Error running TUI: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("👋 Thanks for using Cursor Background Agents CLI!")
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
