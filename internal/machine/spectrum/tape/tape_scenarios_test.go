// Package spectrum implements the ZX Spectrum machine logic.
package tape

import (
	"context"
	"log/slog"
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"os"
	"path/filepath"
	"testing"
)

type tapeLogCaptureHandler struct {
	records []slog.Record
}

func (h *tapeLogCaptureHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *tapeLogCaptureHandler) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}

func (h *tapeLogCaptureHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *tapeLogCaptureHandler) WithGroup(_ string) slog.Handler      { return h }

var tapeScenarios = []dsl.Scenario{
	dsl.NewScenario("Tape block parsing", func(t *testing.T) {
		// Create a fake .tap file
		// Block 1: 3 bytes [0x00, 0x11, 0x22]
		// Block 2: 2 bytes [0xFF, 0xAA]
		data := []byte{
			0x03, 0x00, 0x00, 0x11, 0x22,
			0x02, 0x00, 0xFF, 0xAA,
		}
		filename := filepath.Join(t.TempDir(), "test.tap")
		err := os.WriteFile(filename, data, 0644)
		assert.Equal(t, "WriteFile error", err, nil)

		tape := NewTape()
		assert.Equal(t, "Tape loading info default", tape.LogLoadingInfo, true)
		err = tape.LoadTAP(filename)
		assert.Equal(t, "LoadTAP error", err, nil)
		assert.Equal(t, "Tape filename", tape.Filename, filename)
		assert.Equal(t, "Block count", len(tape.Blocks), 2)
		assert.Equal(t, "Block 1 length", len(tape.Blocks[0]), 3)
		assert.Equal(t, "Block 2 length", len(tape.Blocks[1]), 2)
	}),
	dsl.NewScenario("Tape load logging can be disabled", func(t *testing.T) {
		data := []byte{
			0x03, 0x00, 0x00, 0x11, 0x22,
		}
		filename := filepath.Join(t.TempDir(), "test.tap")
		err := os.WriteFile(filename, data, 0644)
		assert.Equal(t, "WriteFile error", err, nil)

		handler := &tapeLogCaptureHandler{}
		oldDefault := slog.Default()
		slog.SetDefault(slog.New(handler))
		defer slog.SetDefault(oldDefault)

		tape := NewTape()
		tape.LogLoadingInfo = false
		err = tape.LoadTAP(filename)
		assert.Equal(t, "LoadTAP error", err, nil)
		assert.Equal(t, "No tape load log records", len(handler.records), 0)
	}),
	dsl.NewScenario("Tape pulse timing", func(t *testing.T) {
		tape := NewTape()
		tape.Blocks = [][]byte{{0x00, 0x00}} // A dummy block (Flag 0x00 = Header)
		tape.Play()

		// Initial state: Pilot
		assert.Equal(t, "Initial State", tape.State, TapePilot)
		assert.Equal(t, "Initial Pulse Length", tape.PulseLength, uint32(2168))

		// Step less than pulse length
		tape.Step(1000)
		assert.True(t, "Signal remains high", tape.Signal)

		// Step to complete pulse
		tape.Step(1200)
		assert.False(t, "Signal toggles", tape.Signal)
		assert.Equal(t, "Pulse count decremented", tape.PulseCount, 8062)
	}),
}
