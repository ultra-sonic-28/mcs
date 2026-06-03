// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var keyboardScenarios = []dsl.Scenario{
	dsl.NewScenario("Keyboard scanning basic", func(t *testing.T) {
		bus := NewBus()
		// No keys pressed: bits 0-4 should be 1
		// Port 0xFEFE selects Caps Shift row
		val := bus.In(0xFEFE)
		assert.Equal(t, "Scan 0xFEFE (no keys)", val & 0x1F, uint8(0x1F))
		
		// Press 'Z' (Bit 1 of row 0xFE)
		bus.Keyboard.SetKeyState(KeyZ, true)
		val = bus.In(0xFEFE)
		assert.Equal(t, "Scan 0xFEFE (Z pressed)", val & 0x1F, uint8(0x1D)) // 11101
		
		// Release 'Z'
		bus.Keyboard.SetKeyState(KeyZ, false)
		val = bus.In(0xFEFE)
		assert.Equal(t, "Scan 0xFEFE (Z released)", val & 0x1F, uint8(0x1F))
	}),
	dsl.NewScenario("Keyboard multiple rows scanning", func(t *testing.T) {
		bus := NewBus()
		bus.Keyboard.SetKeyState(KeyA, true) // Bit 0 of row 0xFD
		bus.Keyboard.SetKeyState(KeyQ, true) // Bit 0 of row 0xFB
		
		// Scan 0xFDFE (A pressed)
		assert.Equal(t, "Scan 0xFDFE", bus.In(0xFDFE) & 0x1F, uint8(0x1E))
		
		// Scan 0xFBFE (Q pressed)
		assert.Equal(t, "Scan 0xFBFE", bus.In(0xFBFE) & 0x1F, uint8(0x1E))
		
		// Scan with both bits 0: 0xFCFE
		assert.Equal(t, "Scan 0xFCFE (both rows)", bus.In(0xFCFE) & 0x1F, uint8(0x1E))
	}),
	dsl.NewScenario("ULA Out Port 0xFE", func(t *testing.T) {
		bus := NewBus()
		// Border 2 (Red), MIC on, Beeper off
		bus.Out(0x00FE, 0x0A) // 0000 1010 -> Border 2, MIC bit 3=1, Beeper bit 4=0
		assert.Equal(t, "Border Color", bus.BorderColor, uint8(2))
		assert.True(t, "Mic State", bus.MicState)
		assert.False(t, "Beeper State", bus.BeeperState)
		
		// Border 5 (Cyan), Beeper on
		bus.Out(0x00FE, 0x15) // 0001 0101 -> Border 5, Beeper bit 4=1
		assert.Equal(t, "Border Color", bus.BorderColor, uint8(5))
		assert.True(t, "Beeper State", bus.BeeperState)
	}),
	dsl.NewScenario("Tape input bit 6", func(t *testing.T) {
		bus := NewBus()
		bus.TapeInState = true
		assert.Equal(t, "Tape bit 6 high", bus.In(0xFEFE) & 0x40, uint8(0x40))
		
		bus.TapeInState = false
		assert.Equal(t, "Tape bit 6 low", bus.In(0xFEFE) & 0x40, uint8(0x00))
	}),
}
