// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var busScenarios = []dsl.Scenario{
	dsl.NewScenario("ROM is read-only", func(t *testing.T) {
		bus := NewBus()
		// First byte of Spectrum ROM is 0xF3 (DI)
		assert.Equal(t, "Initial ROM byte", bus.Read(0x0000), uint8(0xF3))
		
		bus.Write(0x0000, 0x00)
		assert.Equal(t, "ROM byte after write should remain unchanged", bus.Read(0x0000), uint8(0xF3))
	}),
	dsl.NewScenario("RAM is read-write", func(t *testing.T) {
		bus := NewBus()
		bus.Write(0x4000, 0xAA)
		assert.Equal(t, "RAM byte at 0x4000", bus.Read(0x4000), uint8(0xAA))
		
		bus.Write(0xFFFF, 0x55)
		assert.Equal(t, "RAM byte at 0xFFFF", bus.Read(0xFFFF), uint8(0x55))
	}),
	dsl.NewScenario("Default IO In returns 0xFF", func(t *testing.T) {
		bus := NewBus()
		assert.Equal(t, "IO In at 0x0000", bus.In(0x0000), uint8(0xFF))
		assert.Equal(t, "IO In at 0xFFFF", bus.In(0xFFFF), uint8(0xFF))
	}),
}
