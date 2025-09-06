package cmd

import (
	"fmt"
	"os"

	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/config"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <agent-id>",
	Short: "Get the status and details of a specific background agent",
	Long: `Get the current status and results of a specific background agent.
	
This command shows detailed information about an agent including its status,
source repository, target branch, summary, and creation time.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := config.GetAPIKey()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		agentID := args[0]
		client := client.NewClient(apiKey)

		agent, err := client.GetAgentStatus(agentID)
		if err != nil {
			fmt.Printf("❌ Error getting agent status: %v\n", err)
			os.Exit(1)
		}

		// Display agent information in a nice format
		fmt.Printf("🤖 Agent Details\n")
		fmt.Printf("═══════════════\n\n")

		fmt.Printf("📋 ID: %s\n", agent.ID)
		fmt.Printf("📝 Name: %s\n", agent.Name)
		fmt.Printf("🔄 Status: %s\n", getStatusEmoji(agent.Status))
		fmt.Printf("📅 Created: %s\n", agent.CreatedAt.Format("2006-01-02 15:04:05"))

		fmt.Printf("\n📂 Source Information\n")
		fmt.Printf("─────────────────────\n")
		fmt.Printf("🔗 Repository: %s\n", agent.Source.Repository)
		fmt.Printf("🌿 Reference: %s\n", agent.Source.Ref)

		fmt.Printf("\n🎯 Target Information\n")
		fmt.Printf("─────────────────────\n")
		fmt.Printf("🌿 Branch: %s\n", agent.Target.BranchName)
		fmt.Printf("🔗 Agent URL: %s\n", agent.Target.URL)

		if agent.Target.PrURL != "" {
			fmt.Printf("🔀 Pull Request: %s\n", agent.Target.PrURL)
		}

		fmt.Printf("🔄 Auto Create PR: %t\n", agent.Target.AutoCreatePr)

		if agent.Summary != "" {
			fmt.Printf("\n📄 Summary\n")
			fmt.Printf("──────────\n")
			fmt.Printf("%s\n", agent.Summary)
		}
	},
}

func getStatusEmoji(status string) string {
	switch status {
	case "RUNNING":
		return "🏃 RUNNING"
	case "COMPLETED":
		return "✅ COMPLETED"
	case "FAILED":
		return "❌ FAILED"
	case "CANCELLED":
		return "🚫 CANCELLED"
	default:
		return fmt.Sprintf("❓ %s", status)
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
