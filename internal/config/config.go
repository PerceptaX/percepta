package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Vision  VisionConfig
	Devices map[string]DeviceConfig
}

type VisionConfig struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
}

type DeviceConfig struct {
	Type     string `mapstructure:"type"`
	CameraID string `mapstructure:"camera_id"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "percepta")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set defaults (vision only - no storage in Phase 1)
	viper.SetDefault("vision.provider", "claude")

	// Try to read config (OK if doesn't exist)
	viper.ReadInConfig()

	// Env var overrides
	viper.SetEnvPrefix("PERCEPTA")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override APIKey from env if set
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		cfg.Vision.APIKey = apiKey
	}

	return &cfg, nil
}
