package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the common configuration for all tools
type Config struct {
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
}

// Load loads configuration from a file or environment variables
func Load(appName string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))
	viper.AddConfigPath("/etc/" + appName)

	viper.SetEnvPrefix(strings.ToUpper(appName))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("env", "development")
	viper.SetDefault("log_level", "info")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if we want to rely on defaults/env vars
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
