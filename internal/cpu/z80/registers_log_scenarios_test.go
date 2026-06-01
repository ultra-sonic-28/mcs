package z80

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"mcs/testutils/assert"
	"strings"
	"testing"
)

// logCaptureHandler is a simple slog.Handler that captures log records into a buffer for testing.
type logCaptureHandler struct {
	buf *bytes.Buffer
}

func (h *logCaptureHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *logCaptureHandler) Handle(_ context.Context, r slog.Record) error {
	fmt.Fprintf(h.buf, "[%s] %s", r.Level, r.Message)
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.buf, " %s=%v", a.Key, a.Value)
		return true
	})
	fmt.Fprintln(h.buf)
	return nil
}
func (h *logCaptureHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *logCaptureHandler) WithGroup(name string) slog.Handler        { return h }

// registersLogScenarios defines the test cases for the Registers.LogState method.
var registersLogScenarios = []RegisterScenario{
	{
		Name: "Registers.LogState Output Verification",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.A, r.F = 0x12, 0x34
			r.B, r.C = 0x56, 0x78
			r.D, r.E = 0x9A, 0xBC
			r.H, r.L = 0xDE, 0xF0
			
			r.APrime, r.FPrime = 0x11, 0x22
			r.BPrime, r.CPrime = 0x33, 0x44
			r.DPrime, r.EPrime = 0x55, 0x66
			r.HPrime, r.LPrime = 0x77, 0x88

			r.IX = 0x1234
			r.IY = 0x5678
			r.SP = 0xFFFF
			r.PC = 0x0100
			r.I = 0xAA
			r.R = 0xBB

			var buf bytes.Buffer
			handler := &logCaptureHandler{buf: &buf}
			oldLogger := slog.Default()
			slog.SetDefault(slog.New(handler))
			defer slog.SetDefault(oldLogger)

			r.LogState()

			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			
			assert.Equal(t, "Should have 3 log lines", len(lines), 3)

			// Line 1: Main registers
			assert.Contains(t, "Line 1 should contain 'Register Main'", "Register Main", lines[0])
			assert.Contains(t, "Line 1 should contain A=0x12", "A=0x12", lines[0])
			assert.Contains(t, "Line 1 should contain F=0x34", "F=0x34", lines[0])
			assert.Contains(t, "Line 1 should contain H=0xDE", "H=0xDE", lines[0])

			// Line 2: Alternate registers
			assert.Contains(t, "Line 2 should contain 'Register Alternate'", "Register Alternate", lines[1])
			assert.Contains(t, "Line 2 should contain A'=0x11", "A'=0x11", lines[1])
			assert.Contains(t, "Line 2 should contain H'=0x77", "H'=0x77", lines[1])

			// Line 3: State and control
			assert.Contains(t, "Line 3 should contain 'Register State'", "Register State", lines[2])
			assert.Contains(t, "Line 3 should contain AF=0x1234", "AF=0x1234", lines[2])
			assert.Contains(t, "Line 3 should contain IX=0x1234", "IX=0x1234", lines[2])
			assert.Contains(t, "Line 3 should contain SP=0xFFFF", "SP=0xFFFF", lines[2])
			assert.Contains(t, "Line 3 should contain I=0xAA", "I=0xAA", lines[2])
			assert.Contains(t, "Line 3 should contain R=0xBB", "R=0xBB", lines[2])
		},
	},
}
