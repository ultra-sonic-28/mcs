package spectrum

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var memory48Scenarios = []dsl.Scenario{
	dsl.NewScenario("Memory48 ROM is read-only", func(t *testing.T) {
		mem := NewMemory48()
		// First byte of Spectrum ROM is 0xF3 (DI)
		assert.Equal(t, "Initial ROM byte", mem.Read(0x0000), uint8(0xF3))
		
		mem.Write(0x0000, 0x00)
		assert.Equal(t, "ROM byte after write should remain unchanged", mem.Read(0x0000), uint8(0xF3))
	}),
	dsl.NewScenario("Memory48 RAM is read-write", func(t *testing.T) {
		mem := NewMemory48()
		mem.Write(0x4000, 0xAA)
		assert.Equal(t, "RAM byte at 0x4000", mem.Read(0x4000), uint8(0xAA))
		
		mem.Write(0xFFFF, 0x55)
		assert.Equal(t, "RAM byte at 0xFFFF", mem.Read(0xFFFF), uint8(0x55))
	}),
	dsl.NewScenario("Memory48 Display Memory mapping", func(t *testing.T) {
		mem := NewMemory48()
		mem.Write(0x4000, 0x12)
		displayMem := mem.GetDisplayMemory()
		assert.Equal(t, "First byte of display memory", displayMem[0], uint8(0x12))
	}),
}
