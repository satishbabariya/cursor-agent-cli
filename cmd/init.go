package cmd

import (
	"fmt"
	"os"

	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize cursor-cli with your API key",
	Long: `Initialize cursor-cli by setting up your Cursor API key.

You can get your API key from the Cursor Dashboard → Integrations.
The API key will be stored in ~/.cursor-cli.yaml for future use.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Welcome to cursor-cli setup!")
		fmt.Println()
		fmt.Println("To get started, you'll need a Cursor API key.")
		fmt.Println("You can create one at: https://cursor.com/dashboard")
		fmt.Println("Navigate to: Dashboard → Integrations → Create API Key")
		fmt.Println()

		fmt.Print("Please enter your Cursor API key: ")
		var apiKey string
		fmt.Scanln(&apiKey)

		if apiKey == "" {
			fmt.Println("❌ Error: API key cannot be empty")
			os.Exit(1)
		}

		// Test the API key by making a request to the API key info endpoint
		fmt.Println("🔍 Validating API key...")
		client := client.NewClient(apiKey)
		keyInfo, err := client.GetAPIKeyInfo()
		if err != nil {
			fmt.Printf("❌ Error: Invalid API key or network error: %v\n", err)
			os.Exit(1)
		}

		// Save the API key to config
		if err := config.SaveAPIKey(apiKey); err != nil {
			fmt.Printf("❌ Error saving API key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ API key validated and saved successfully!")
		fmt.Printf("📧 Authenticated as: %s\n", keyInfo.UserEmail)
		fmt.Printf("🔑 Key ID: %s\n", keyInfo.ID)
		fmt.Printf("📅 Created: %s\n", keyInfo.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
		fmt.Println("🎉 cursor-cli is ready to use!")
		fmt.Println("Try running: cursor-cli list")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
