// Package config provides configuration management for the MCS project.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// AcceptedLoggingLevels lists the logging levels accepted in config.json.
var AcceptedLoggingLevels = []string{"INFO", "DEBUG", "WARN", "ERROR"}

// AcceptedScaleValues lists the display scale values accepted in config.json.
var AcceptedScaleValues = []int{1, 2, 3, 4}

// BadLoggingLevelError reports an unsupported logging level from config.json.
type BadLoggingLevelError struct {
	Level          string
	AcceptedValues []string
}

func (e *BadLoggingLevelError) Error() string {
	return fmt.Sprintf("bad logging level %q; accepted values: %v", e.Level, e.AcceptedValues)
}

// BadScaleError reports an unsupported display scale from config.json.
type BadScaleError struct {
	Scale          int
	AcceptedValues []int
}

func (e *BadScaleError) Error() string {
	return fmt.Sprintf("bad scale %d; accepted values: %v", e.Scale, e.AcceptedValues)
}

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

// LoggingConfig holds configuration for application logging.
type LoggingConfig struct {
	// Enabled controls whether logging is written.
	Enabled bool `json:"enabled"`
	// Level is the configured logging level (DEBUG, INFO, WARN, ERROR).
	Level string `json:"level"`
	// Z80 defines logging settings for the Z80 CPU.
	Z80 Z80LoggingConfig `json:"z80"`
}

// Z80LoggingConfig holds Z80-specific logging configuration.
type Z80LoggingConfig struct {
	// Instructions controls whether all registered Z80 instructions are logged at startup.
	Instructions bool `json:"instructions"`
	// Tape controls whether tape loading information is logged.
	Tape bool `json:"tape"`
}

// Config represents the application configuration.
type Config struct {
	Logging LoggingConfig `json:"logging"`
	Display DisplayConfig `json:"display"`
}

// Load reads the configuration from the specified file path.
// If the file does not exist, it creates it with default values.
// It always returns a valid configuration.
func Load(filePath string) (*Config, error) {
	defaultCfg := &Config{
		Logging: LoggingConfig{
			Enabled: false,
			Level:   "INFO",
			Z80: Z80LoggingConfig{
				Instructions: true,
				Tape:         true,
			},
		},
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

	// Try to read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create it with defaults
			data, err := json.MarshalIndent(defaultCfg, "", "  ")
			if err == nil {
				_ = os.WriteFile(filePath, data, 0644)
			}
			return defaultCfg, nil
		}
		// Other reading error, return defaults
		return defaultCfg, nil
	}

	cfg := *defaultCfg
	if err := json.Unmarshal(data, &cfg); err != nil {
		// Decoding error, return defaults
		return defaultCfg, nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err == nil {
		if displayRaw, ok := raw["display"]; ok {
			var displayMap map[string]json.RawMessage
			if err := json.Unmarshal(displayRaw, &displayMap); err == nil {
				if _, hasScale := displayMap["scale"]; !hasScale {
					cfg.Display.Scale = defaultCfg.Display.Scale
				}
			}
		}
	}

	level := strings.ToUpper(cfg.Logging.Level)
	if !isAcceptedLoggingLevel(level) {
		return nil, &BadLoggingLevelError{
			Level:          cfg.Logging.Level,
			AcceptedValues: AcceptedLoggingLevels,
		}
	}
	cfg.Logging.Level = level

	if !isAcceptedScale(cfg.Display.Scale) {
		return nil, &BadScaleError{
			Scale:          cfg.Display.Scale,
			AcceptedValues: AcceptedScaleValues,
		}
	}

	return &cfg, nil
}

func isAcceptedLoggingLevel(level string) bool {
	for _, acceptedLevel := range AcceptedLoggingLevels {
		if level == acceptedLevel {
			return true
		}
	}
	return false
}

func isAcceptedScale(scale int) bool {
	for _, acceptedScale := range AcceptedScaleValues {
		if scale == acceptedScale {
			return true
		}
	}
	return false
}
