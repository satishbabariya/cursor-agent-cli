package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cursor-cli",
	Short: "CLI tool for managing Cursor Background Agents",
	Long: `A command-line interface application to interact with Cursor's Background Agents API.
	
This tool allows you to:
- List all background agents
- Get agent status and details
- View agent conversations
- Add follow-up instructions to agents
- Manage API keys and configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to cursor-cli! Use --help to see available commands.")
		fmt.Println("Run 'cursor-cli init' to set up your API key.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cursor-cli.yaml)")
	rootCmd.PersistentFlags().String("api-key", "", "Cursor API key (can also be set via CURSOR_API_KEY env var)")

	// Bind the flag to viper
	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cursor-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cursor-cli")
	}

	viper.AutomaticEnv()         // read in environment variables that match
	viper.SetEnvPrefix("CURSOR") // will be uppercased automatically

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
