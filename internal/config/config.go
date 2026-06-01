// Package config provides configuration management for the MCS project.
package config

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration.
type Config struct {
	LoggingEnabled bool   `json:"logging_enabled"`
	LogLevel       string `json:"log_level"`
}

// Load reads the configuration from the specified file path.
// If the file does not exist, it creates it with default values.
// It always returns a valid configuration.
func Load(filePath string) (*Config, error) {
	defaultCfg := &Config{
		LoggingEnabled: false,
		LogLevel:       "INFO",
	}

	// Try to open the file
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create it with defaults
			data, err := json.MarshalIndent(defaultCfg, "", "  ")
			if err == nil {
				_ = os.WriteFile(filePath, data, 0644)
			}
			return defaultCfg, nil
		}
		// Other opening error, return defaults
		return defaultCfg, nil
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		// Decoding error, return defaults
		return defaultCfg, nil
	}

	return &cfg, nil
}
