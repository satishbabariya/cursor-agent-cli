package cmd

import (
	"fmt"
	"os"

	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/config"
	"github.com/spf13/cobra"
)

// followupCmd represents the followup command
var followupCmd = &cobra.Command{
	Use:   "followup <agent-id> <prompt>",
	Short: "Add a follow-up instruction to a running background agent",
	Long: `Send an additional instruction to a running background agent.
	
This allows you to provide additional context or modify the agent's task
while it's still running. The prompt will be added to the agent's conversation.

Example:
  cursor-cli followup bc_abc123 "Also add a section about troubleshooting"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := config.GetAPIKey()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		agentID := args[0]
		prompt := args[1]

		client := client.NewClient(apiKey)

		fmt.Printf("📤 Sending follow-up instruction to agent %s...\n", agentID)

		response, err := client.AddFollowup(agentID, prompt)
		if err != nil {
			fmt.Printf("❌ Error adding follow-up: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Follow-up instruction sent successfully!\n")
		fmt.Printf("🤖 Agent ID: %s\n", response.ID)
		fmt.Printf("💬 Instruction: %s\n", prompt)
		fmt.Println()
		fmt.Println("The agent will process your follow-up instruction and continue working.")
		fmt.Printf("You can check the status with: cursor-cli status %s\n", agentID)
	},
}

func init() {
	rootCmd.AddCommand(followupCmd)
}
