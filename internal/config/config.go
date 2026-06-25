// Package config provides configuration management for the MCS project.
package config

import (
	"encoding/json"
	"os"
)

// BorderConfig holds configuration for the CRT-like display border.
type BorderConfig struct {
	// Color is the hexadecimal color code (e.g., "#D6CDC9") of the border.
	Color string `json:"color"`
	// Width is the thickness of the border in logical pixels.
	Width int `json:"width"`
}

// ToolbarConfig holds configuration for the top toolbar.
type ToolbarConfig struct {
	// Color is the hexadecimal color code (e.g., "#D6CDC9") of the toolbar.
	Color string `json:"color"`
	// Height is the height of the toolbar in logical pixels.
	Height int `json:"height"`
}

// DisplayConfig holds configuration for the emulator display settings.
type DisplayConfig struct {
	// Scale is the scaling factor for the display.
	Scale int `json:"scale"`
	// Border defines the settings for the CRT-like screen border.
	Border BorderConfig `json:"border"`
	// Toolbar defines the settings for the top toolbar.
	Toolbar ToolbarConfig `json:"toolbar"`
}

// Config represents the application configuration.
type Config struct {
	LoggingEnabled bool          `json:"logging_enabled"`
	LogLevel       string        `json:"log_level"`
	Display        DisplayConfig `json:"display"`
}

// Load reads the configuration from the specified file path.
// If the file does not exist, it creates it with default values.
// It always returns a valid configuration.
func Load(filePath string) (*Config, error) {
	defaultCfg := &Config{
		LoggingEnabled: false,
		LogLevel:       "INFO",
		Display: DisplayConfig{
			Border: BorderConfig{
				Color: "#D6EFC9",
				Width: 15,
			},
			Toolbar: ToolbarConfig{
				Color:  "#D6CDC9",
				Height: 20,
			},
			Scale: 2,
		},
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
