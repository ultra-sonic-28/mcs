// Package logger provides a customized slog infrastructure for the MCS project.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// LogHandler implements a custom slog.Handler that matches the project's requirements.
// Format: YYYY-MM-DD HH:MM:SS,ms [LEVEL] AppName: Message, key1=val1, key2=val2, ...
type LogHandler struct {
	w       io.Writer
	level   slog.Level
	enabled bool
}

// Enabled reports whether the handler handles records at the given level.
func (h *LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	if !h.enabled {
		return false
	}
	return level >= h.level
}

// Handle handles the Record.
func (h *LogHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String()
	timestamp := r.Time.Format("2006-01-02 15:04:05,000")
	msg := fmt.Sprintf("%s [%s] MCS: %s", timestamp, level, r.Message)

	r.Attrs(func(a slog.Attr) bool {
		msg += fmt.Sprintf(", %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintln(h.w, msg)
	return nil
}

// WithAttrs returns a new Handler with the given attributes.
func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new Handler with the given group name.
func (h *LogHandler) WithGroup(name string) slog.Handler {
	return h
}

// Setup initializes the logging system.
// If enabled is false, it removes the log file if it exists and disables all output.
// If enabled is true, it opens the log file (truncating it), sets up a MultiWriter
// to log to both the file and stdout, and sets the default slog logger.
// It returns a cleanup function to close the log file and any error encountered.
func Setup(filePath string, enabled bool, levelStr string) (func() error, error) {
	var level slog.Level
	switch levelStr {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var logFile *os.File
	var mw io.Writer = os.Stdout

	if enabled {
		var err error
		logFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		mw = io.MultiWriter(os.Stdout, logFile)
	} else {
		// Remove log file if it exists
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			// Log error to stderr since logging is disabled
			fmt.Fprintf(os.Stderr, "failed to remove log file: %v\n", err)
		}
		mw = io.Discard
	}

	handler := &LogHandler{w: mw, level: level, enabled: enabled}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	cleanup := func() error {
		if logFile != nil {
			return logFile.Close()
		}
		return nil
	}

	return cleanup, nil
}
