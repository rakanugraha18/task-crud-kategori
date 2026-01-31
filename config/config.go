package config

import "github.com/spf13/viper"

// Config holds the application configuration
type Config struct {
	AppPort string
	DBPath  string
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Check if .env file exists
	viper.SetConfigFile(".env")
	// Set the file type to env
	viper.SetConfigType("env")
	// Load environment variables
	viper.AutomaticEnv()
	// Read .env file if it exists
	err := viper.ReadInConfig()
	// Handle error
	if err != nil {
		return nil, err
	}
	// Map configuration to struct
	return &Config{
		AppPort: viper.GetString("APP_PORT"),
		DBPath:  viper.GetString("DB_CONN"),
	}, nil
}
