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
		assert.Equal(t, "Logging Enabled default", cfg.Logging.Enabled, false)
		assert.Equal(t, "Logging Level default", cfg.Logging.Level, "INFO")
		assert.Equal(t, "Z80 instructions logging default", cfg.Logging.Z80.Instructions, true)
		assert.Equal(t, "Z80 tape logging default", cfg.Logging.Z80.Tape, true)
		assert.Equal(t, "Border Color default", cfg.Display.Border.Color, "#D6EFC9")
		assert.Equal(t, "Border Width default", cfg.Display.Border.Width, 15)
		assert.Equal(t, "Toolbar Color default", cfg.Display.Toolbar.Color, "#D6CDC9")
		assert.Equal(t, "Toolbar Height default", cfg.Display.Toolbar.Height, 20)

		// Verify the file was created on disk
		_, err = os.Stat(configPath)
		assert.Equal(t, "Config file should exist", err, nil)
	}),

	dsl.NewScenario("Load configured settings from existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		expectedCfg := &Config{
			Logging: LoggingConfig{
				Enabled: true,
				Level:   "DEBUG",
				Z80: Z80LoggingConfig{
					Instructions: false,
					Tape:         false,
				},
			},
			Display: DisplayConfig{
				Border: BorderConfig{
					Color: "#FF0000",
					Width: 30,
				},
				Toolbar: ToolbarConfig{
					Color:  "#00FF00",
					Height: 15,
				},
			},
		}

		data, err := json.MarshalIndent(expectedCfg, "", "  ")
		assert.Equal(t, "Marshal error should be nil", err, nil)

		err = os.WriteFile(configPath, data, 0644)
		assert.Equal(t, "WriteFile error should be nil", err, nil)

		cfg, err := Load(configPath)
		assert.Equal(t, "Load error should be nil", err, nil)
		assert.Equal(t, "Logging Enabled", cfg.Logging.Enabled, true)
		assert.Equal(t, "Logging Level", cfg.Logging.Level, "DEBUG")
		assert.Equal(t, "Z80 instructions logging", cfg.Logging.Z80.Instructions, false)
		assert.Equal(t, "Z80 tape logging", cfg.Logging.Z80.Tape, false)
		assert.Equal(t, "Border Color", cfg.Display.Border.Color, "#FF0000")
		assert.Equal(t, "Border Width", cfg.Display.Border.Width, 30)
		assert.Equal(t, "Toolbar Color", cfg.Display.Toolbar.Color, "#00FF00")
		assert.Equal(t, "Toolbar Height", cfg.Display.Toolbar.Height, 15)
	}),

	dsl.NewScenario("Normalize logging level casing", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		data := []byte(`{
  "logging": {
    "enabled": true,
    "level": "wArN"
  }
}`)

		err := os.WriteFile(configPath, data, 0644)
		assert.Equal(t, "WriteFile error should be nil", err, nil)

		cfg, err := Load(configPath)
		assert.Equal(t, "Load error should be nil", err, nil)
		assert.Equal(t, "Logging Level normalized", cfg.Logging.Level, "WARN")
	}),

	dsl.NewScenario("Reject invalid logging level", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		data := []byte(`{
  "logging": {
    "enabled": true,
    "level": "verbose"
  }
}`)

		err := os.WriteFile(configPath, data, 0644)
		assert.Equal(t, "WriteFile error should be nil", err, nil)

		cfg, err := Load(configPath)
		assert.Equal(t, "Config should be nil", cfg, (*Config)(nil))
		assert.True(t, "Load error should be BadLoggingLevelError", err != nil)

		badLoggingLevel, ok := err.(*BadLoggingLevelError)
		assert.True(t, "Error should be BadLoggingLevelError", ok)
		assert.Equal(t, "Bad logging level", badLoggingLevel.Level, "verbose")
		assert.DeepEqual(t, "Accepted logging levels", badLoggingLevel.AcceptedValues, AcceptedLoggingLevels)
	}),

	dsl.NewScenario("Load default Z80 logging settings when existing file omits them", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		data := []byte(`{
  "logging": {
    "enabled": true,
    "level": "INFO"
  }
}`)

		err := os.WriteFile(configPath, data, 0644)
		assert.Equal(t, "WriteFile error should be nil", err, nil)

		cfg, err := Load(configPath)
		assert.Equal(t, "Load error should be nil", err, nil)
		assert.Equal(t, "Z80 instructions logging default", cfg.Logging.Z80.Instructions, true)
		assert.Equal(t, "Z80 tape logging default", cfg.Logging.Z80.Tape, true)
	}),
}
