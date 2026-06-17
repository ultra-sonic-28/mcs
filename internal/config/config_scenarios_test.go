// Package config defines the scenarios for config tests.
package config

import (
	"encoding/json"
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"os"
	"path/filepath"
	"testing"
)

// configScenarios holds the list of scenarios to run.
var configScenarios = []dsl.Scenario{
	dsl.NewScenario("Load default config when file does not exist", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		cfg, err := Load(configPath)
		assert.Equal(t, "Load error should be nil", err, nil)
		assert.Equal(t, "LoggingEnabled default", cfg.LoggingEnabled, false)
		assert.Equal(t, "LogLevel default", cfg.LogLevel, "INFO")
		assert.Equal(t, "Border Color default", cfg.Display.Border.Color, "#D6CDC9")
		assert.Equal(t, "Border Width default", cfg.Display.Border.Width, 0)

		// Verify the file was created on disk
		_, err = os.Stat(configPath)
		assert.Equal(t, "Config file should exist", err, nil)
	}),

	dsl.NewScenario("Load configured settings from existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		expectedCfg := &Config{
			LoggingEnabled: true,
			LogLevel:       "DEBUG",
			Display: DisplayConfig{
				Border: BorderConfig{
					Color: "#FF0000",
					Width: 30,
				},
			},
		}

		data, err := json.MarshalIndent(expectedCfg, "", "  ")
		assert.Equal(t, "Marshal error should be nil", err, nil)

		err = os.WriteFile(configPath, data, 0644)
		assert.Equal(t, "WriteFile error should be nil", err, nil)

		cfg, err := Load(configPath)
		assert.Equal(t, "Load error should be nil", err, nil)
		assert.Equal(t, "LoggingEnabled", cfg.LoggingEnabled, true)
		assert.Equal(t, "LogLevel", cfg.LogLevel, "DEBUG")
		assert.Equal(t, "Border Color", cfg.Display.Border.Color, "#FF0000")
		assert.Equal(t, "Border Width", cfg.Display.Border.Width, 30)
	}),
}
