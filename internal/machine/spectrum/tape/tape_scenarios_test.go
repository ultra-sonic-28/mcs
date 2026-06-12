// Package spectrum implements the ZX Spectrum machine logic.
package tape

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"os"
	"testing"
)

var tapeScenarios = []dsl.Scenario{
	dsl.NewScenario("Tape block parsing", func(t *testing.T) {
		// Create a fake .tap file
		// Block 1: 3 bytes [0x00, 0x11, 0x22]
		// Block 2: 2 bytes [0xFF, 0xAA]
		data := []byte{
			0x03, 0x00, 0x00, 0x11, 0x22,
			0x02, 0x00, 0xFF, 0xAA,
		}
		filename := "test.tap"
		os.WriteFile(filename, data, 0644)
		defer os.Remove(filename)
		
		tape := NewTape()
		err := tape.LoadTAP(filename)
		assert.Equal(t, "LoadTAP error", err, nil)
		assert.Equal(t, "Block count", len(tape.Blocks), 2)
		assert.Equal(t, "Block 1 length", len(tape.Blocks[0]), 3)
		assert.Equal(t, "Block 2 length", len(tape.Blocks[1]), 2)
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
