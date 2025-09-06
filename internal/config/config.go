package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

var (
	ErrAPIKeyNotSet = errors.New("API key not set. Please run 'cursor-cli init' or set CURSOR_API_KEY environment variable")
)

// GetAPIKey returns the API key from config or environment
func GetAPIKey() (string, error) {
	// First check if it's set via flag/config
	apiKey := viper.GetString("api_key")
	if apiKey != "" {
		return apiKey, nil
	}

	// Then check environment variable
	apiKey = os.Getenv("CURSOR_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}

	return "", ErrAPIKeyNotSet
}

// SaveAPIKey saves the API key to the config file
func SaveAPIKey(apiKey string) error {
	viper.Set("api_key", apiKey)

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Set config file path
	configFile := home + "/.cursor-cli.yaml"

	return viper.WriteConfigAs(configFile)
}
