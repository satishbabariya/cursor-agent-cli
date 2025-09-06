package cmd

import (
	"fmt"
	"os"

	"github.com/satishbabariya/cursor-background-agent-cli/internal/client"
	"github.com/satishbabariya/cursor-background-agent-cli/internal/config"
	"github.com/spf13/cobra"
)

// keyinfoCmd represents the keyinfo command
var keyinfoCmd = &cobra.Command{
	Use:   "keyinfo",
	Short: "Display information about your API key",
	Long: `Retrieve and display information about your current API key.
	
This command shows details about the authenticated API key including
the key ID, name, creation date, and associated user email.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := config.GetAPIKey()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		client := client.NewClient(apiKey)

		keyInfo, err := client.GetAPIKeyInfo()
		if err != nil {
			fmt.Printf("❌ Error getting API key info: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🔑 API Key Information\n")
		fmt.Printf("═════════════════════\n\n")

		fmt.Printf("📋 ID: %s\n", keyInfo.ID)
		fmt.Printf("📝 Name: %s\n", keyInfo.Name)
		fmt.Printf("📧 User Email: %s\n", keyInfo.UserEmail)
		fmt.Printf("📅 Created: %s\n", keyInfo.CreatedAt.Format("2006-01-02 15:04:05"))

		fmt.Println("\n✅ API key is valid and active!")
	},
}

func init() {
	rootCmd.AddCommand(keyinfoCmd)
}
