package cmd

import (
	"fmt"
	"os"

	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/config"
	"github.com/spf13/cobra"
)

// conversationCmd represents the conversation command
var conversationCmd = &cobra.Command{
	Use:   "conversation <agent-id>",
	Short: "Get the conversation history of a background agent",
	Long: `Retrieve the conversation history of a background agent.
	
This command shows all messages in the agent's conversation, including
user messages and agent responses.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := config.GetAPIKey()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		agentID := args[0]
		client := client.NewClient(apiKey)

		conversation, err := client.GetAgentConversation(agentID)
		if err != nil {
			fmt.Printf("❌ Error getting agent conversation: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("💬 Conversation History for Agent: %s\n", conversation.ID)
		fmt.Printf("═══════════════════════════════════════\n\n")

		if len(conversation.Messages) == 0 {
			fmt.Println("📭 No messages found in this conversation.")
			return
		}

		for i, message := range conversation.Messages {
			emoji := getMessageTypeEmoji(message.Type)
			fmt.Printf("%s Message %d (ID: %s)\n", emoji, i+1, message.ID)
			fmt.Printf("─────────────────────────────────\n")
			fmt.Printf("%s\n", message.Text)

			if i < len(conversation.Messages)-1 {
				fmt.Println()
			}
		}
	},
}

func getMessageTypeEmoji(messageType string) string {
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

func init() {
	rootCmd.AddCommand(conversationCmd)
}
