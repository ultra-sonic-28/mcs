package logger

import (
	"bytes"
	"context"
	"log/slog"
	"mcs/testutils"
	"mcs/testutils/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type LoggerScenario struct {
	Name string
	Run  func(t *testing.T)
}

var loggerScenarios = []LoggerScenario{
	{
		Name: "LogHandler Formatting",
		Run: func(t *testing.T) {
			var buf bytes.Buffer
			handler := &LogHandler{w: &buf, enabled: true}
			logger := slog.New(handler)

			logger.Info("test message", "key", "value")

			output := buf.String()
			// Expected format: YYYY-MM-DD HH:MM:SS,ms [INFO] MCS: test message, key=value
			assert.Contains(t, "Should contain level", " [INFO] ", output)
			assert.Contains(t, "Should contain app name", " MCS: ", output)
			assert.Contains(t, "Should contain message", " test message", output)
			assert.Contains(t, "Should contain attribute", "key=value", output)

			// Verify timestamp format (basic check)
			parts := strings.Split(output, " ")
			assert.True(t, "Timestamp should be present", len(parts) >= 2)
		},
	},
	{
		Name: "LogHandler Enabled",
		Run: func(t *testing.T) {
			handler := &LogHandler{level: slog.LevelInfo, enabled: true}
			assert.True(t, "Should be enabled for Info", handler.Enabled(context.Background(), slog.LevelInfo))
			assert.False(t, "Should NOT be enabled for Debug when level is Info", handler.Enabled(context.Background(), slog.LevelDebug))
			assert.True(t, "Should be enabled for Error", handler.Enabled(context.Background(), slog.LevelError))

			handlerDebug := &LogHandler{level: slog.LevelDebug, enabled: true}
			assert.True(t, "Should be enabled for Debug when level is Debug", handlerDebug.Enabled(context.Background(), slog.LevelDebug))

			handlerDisabled := &LogHandler{level: slog.LevelDebug, enabled: false}
			assert.False(t, "Should be disabled when enabled is false", handlerDisabled.Enabled(context.Background(), slog.LevelError))
		},
	},
	{
		Name: "LogHandler WithAttrs and WithGroup",
		Run: func(t *testing.T) {
			handler := &LogHandler{}
			h2 := handler.WithAttrs([]slog.Attr{slog.String("k", "v")})
			assert.Equal(t, "WithAttrs should return same handler", h2, slog.Handler(handler))

			h3 := handler.WithGroup("test")
			assert.Equal(t, "WithGroup should return same handler", h3, slog.Handler(handler))
		},
	},
	{
		Name: "LogHandler Multiple Attributes",
		Run: func(t *testing.T) {
			var buf bytes.Buffer
			handler := &LogHandler{w: &buf, enabled: true}
			logger := slog.New(handler)

			logger.Error("error message", "err", "failed", "code", 500)

			output := buf.String()
			assert.Contains(t, "Should contain error level", " [ERROR] ", output)
			assert.Contains(t, "Should contain message", " error message", output)
			assert.Contains(t, "Should contain first attribute", "err=failed", output)
			assert.Contains(t, "Should contain second attribute", "code=500", output)
		},
	},
	{
		Name: "Logger Setup",
		Run: func(t *testing.T) {
			tempDir := t.TempDir()
			logPath := filepath.Join(tempDir, "test.log")

			oldLogger := slog.Default()
			defer slog.SetDefault(oldLogger)

			var cleanup func() error
			var err error

			stdoutOutput := testutils.CaptureStdout(t, func() {
				cleanup, err = Setup(logPath, true, "INFO")
				if err != nil {
					return
				}
				slog.Info("setup test message")
				slog.Debug("should not appear")
			})

			assert.True(t, "Setup should not return error", err == nil)
			assert.True(t, "Cleanup should not be nil", cleanup != nil)

			if cleanup != nil {
				err = cleanup()
				assert.True(t, "Cleanup should not return error", err == nil)
			}

			// Verify stdout
			assert.Contains(t, "Stdout should contain message", "setup test message", stdoutOutput)
			assert.False(t, "Stdout should NOT contain debug message", strings.Contains(stdoutOutput, "should not appear"))

			// Verify file
			content, err := os.ReadFile(logPath)
			assert.True(t, "Should be able to read log file", err == nil)
			assert.Contains(t, "Log file should contain message", "setup test message", string(content))
			assert.False(t, "Log file should NOT contain debug message", strings.Contains(string(content), "should not appear"))
		},
	},
	{
		Name: "Logger Setup Disabled",
		Run: func(t *testing.T) {
			tempDir := t.TempDir()
			logPath := filepath.Join(tempDir, "disabled.log")

			stdoutOutput := testutils.CaptureStdout(t, func() {
				cleanup, err := Setup(logPath, false, "INFO")
				assert.True(t, "Setup should not return error", err == nil)
				if cleanup != nil {
					defer cleanup()
				}
				slog.Info("stdout only message")
			})

			assert.False(t, "Stdout should NOT contain message when logging is disabled", strings.Contains(stdoutOutput, "stdout only message"))

			// Verify file does not exist
			_, err := os.Stat(logPath)
			assert.True(t, "Log file should not exist", os.IsNotExist(err))
		},
	},
	{
		Name: "Logger Setup Truncate",
		Run: func(t *testing.T) {
			tempDir := t.TempDir()
			logPath := filepath.Join(tempDir, "truncate.log")

			// Pre-fill file
			err := os.WriteFile(logPath, []byte("old content"), 0666)
			assert.True(t, "Should be able to write pre-fill", err == nil)

			cleanup, err := Setup(logPath, true, "INFO")
			assert.True(t, "Setup should not return error", err == nil)
			if cleanup != nil {
				defer cleanup()
			}

			content, err := os.ReadFile(logPath)
			assert.True(t, "Should be able to read log file", err == nil)
			assert.Equal(t, "Log file should be truncated", string(content), "")
		},
	},
	{
		Name: "Logger Setup Error",
		Run: func(t *testing.T) {
			// Use an invalid path. On Windows, a path with invalid characters should fail.
			_, err := Setup(`Z:\non\existent\path\to\log.log`, true, "INFO")
			assert.True(t, "Setup should return error for invalid path", err != nil)
		},
	},
}
