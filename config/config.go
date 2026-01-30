package config

import "github.com/spf13/viper"

type Config struct {
	AppPort string
	DBPath  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		AppPort: viper.GetString("APP_PORT"),
		DBPath:  viper.GetString("DB_CONN"),
	}, nil
}
