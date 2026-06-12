package spectrum

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var memory128Scenarios = []dsl.Scenario{
	dsl.NewScenario("Memory128 RAM paging", func(t *testing.T) {
		mem := NewMemory128()
		
		// Default: RAM 0 at 0xC000
		mem.Write(0xC000, 0x42)
		assert.Equal(t, "Value in RAM 0", mem.Read(0xC000), uint8(0x42))
		
		// Page in RAM 1 at 0xC000 (val 1)
		mem.Page(0x01)
		assert.Equal(t, "Value in RAM 1 (should be 0)", mem.Read(0xC000), uint8(0x00))
		mem.Write(0xC000, 0x43)
		assert.Equal(t, "Value in RAM 1", mem.Read(0xC000), uint8(0x43))
		
		// Page back RAM 0
		mem.Page(0x00)
		assert.Equal(t, "Value in RAM 0 again", mem.Read(0xC000), uint8(0x42))
	}),
	dsl.NewScenario("Memory128 ROM paging", func(t *testing.T) {
		mem := NewMemory128()
		
		// Default: ROM 0
		assert.Equal(t, "Initial ROM 0 byte", mem.Read(0x0000), uint8(0xF3))
		assert.True(t, "ROM 1 should NOT be active", !mem.IsRom1Active())
		
		// Page in ROM 1 (bit 4 = 1)
		mem.Page(0x10)
		assert.Equal(t, "ROM 1 byte", mem.Read(0x0000), uint8(0xF3))
		assert.True(t, "ROM 1 should be active", mem.IsRom1Active())
	}),
	dsl.NewScenario("Memory128 Paging Lock", func(t *testing.T) {
		mem := NewMemory128()
		
		// Lock paging (bit 5 = 1)
		mem.Page(0x20)
		
		// Try to page in RAM 1
		mem.Page(0x01)
		
		// Should still be RAM 0
		mem.Write(0xC000, 0xAA)
		mem.Page(0x00) // Attempt to page back (should be ignored)
		assert.Equal(t, "RAM should still be locked", mem.Read(0xC000), uint8(0xAA))
	}),
	dsl.NewScenario("Memory128 Shadow Screen paging", func(t *testing.T) {
		mem := NewMemory128()
		
		// Default screen is 5
		mem.Write(0x4000, 0x55) // RAM 5 starts at 0x4000
		displayMem := mem.GetDisplayMemory()
		assert.Equal(t, "Screen 5 first byte", displayMem[0], uint8(0x55))
		
		// Switch to screen 7 (bit 3 = 1)
		mem.Page(0x08)
		mem.Write(0xC000, 0x77) // Map RAM 7 at 0xC000 to write to it easily
		// Wait, 0xC000 is for currentRam. 
		// If currentRam is 0 (default), 0xC000 is RAM 0.
		// If we want to write to RAM 7, we should page it in.
		mem.Page(0x08 | 0x07) // Screen 7, RAM 7 at 0xC000
		mem.Write(0xC000, 0x77)
		
		displayMem = mem.GetDisplayMemory()
		assert.Equal(t, "Screen 7 first byte", displayMem[0], uint8(0x77))
	}),
}
