package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all background agents",
	Long: `Retrieve and display a list of all background agents associated with your account.
	
This command shows the ID, name, status, and creation time of each agent.
By default, expired agents are filtered out and only running/finished agents are shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := config.GetAPIKey()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		showAll, _ := cmd.Flags().GetBool("all")

		client := client.NewClient(apiKey)
		response, err := client.ListAgents(limit, cursor)
		if err != nil {
			fmt.Printf("❌ Error listing agents: %v\n", err)
			os.Exit(1)
		}

		// Filter agents if not showing all
		filteredAgents := response.Agents
		if !showAll {
			filteredAgents = filterActiveAgents(response.Agents)
		}

		if len(filteredAgents) == 0 {
			if showAll {
				fmt.Println("📭 No background agents found.")
			} else {
				fmt.Println("📭 No active background agents found.")
				fmt.Println("💡 Use --all flag to see expired agents as well.")
			}
			return
		}

		statusText := "all"
		if !showAll {
			statusText = "active"
		}
		fmt.Printf("📋 Found %d %s background agents:\n\n", len(filteredAgents), statusText)

		// Create a tab writer for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tREPOSITORY\tCREATED")
		fmt.Fprintln(w, "──\t────\t──────\t──────────\t───────")

		for _, agent := range filteredAgents {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				agent.ID,
				agent.Name,
				agent.Status,
				agent.Source.Repository,
				agent.CreatedAt.Format("2006-01-02 15:04"),
			)
		}

		w.Flush()

		if response.NextCursor != "" {
			fmt.Printf("\n🔗 More results available. Use --cursor=%s to get the next page.\n", response.NextCursor)
		}
	},
}

// filterActiveAgents filters out expired agents, keeping only running/finished ones
func filterActiveAgents(agents []client.Agent) []client.Agent {
	var filtered []client.Agent
	for _, agent := range agents {
		status := strings.ToUpper(agent.Status)
		// Include running, completed, failed, cancelled - exclude expired
		if status != "EXPIRED" {
			filtered = append(filtered, agent)
		}
	}
	return filtered
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().IntP("limit", "l", 20, "Number of agents to return (1-100)")
	listCmd.Flags().StringP("cursor", "c", "", "Pagination cursor from previous response")
	listCmd.Flags().BoolP("all", "a", false, "Show all agents including expired ones")
}
