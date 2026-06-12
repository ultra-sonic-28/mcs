package spectrum

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var bus128Scenarios = []dsl.Scenario{
	dsl.NewScenario("Spectrum 128K RAM paging", func(t *testing.T) {
		bus := NewBus128()
		
		// Default: ROM 0, RAM 0 at 0xC000
		bus.Write(0xC000, 0x42)
		assert.Equal(t, "Value in RAM 0", bus.Read(0xC000), uint8(0x42))
		
		// Page in RAM 1 at 0xC000 (Port 0x7FFD, val 1)
		bus.Out(0x7FFD, 0x01)
		assert.Equal(t, "Value in RAM 1 (should be 0)", bus.Read(0xC000), uint8(0x00))
		bus.Write(0xC000, 0x43)
		assert.Equal(t, "Value in RAM 1", bus.Read(0xC000), uint8(0x43))
		
		// Page back RAM 0
		bus.Out(0x7FFD, 0x00)
		assert.Equal(t, "Value in RAM 0 again", bus.Read(0xC000), uint8(0x42))
	}),
	dsl.NewScenario("Spectrum 128K ROM paging", func(t *testing.T) {
		bus := NewBus128()
		
		// Default: ROM 0
		// First byte of 128K ROM 0 is 0xF3
		assert.Equal(t, "Initial ROM 0 byte", bus.Read(0x0000), uint8(0xF3))
		
		// Page in ROM 1 (Port 0x7FFD, bit 4 = 1)
		bus.Out(0x7FFD, 0x10)
		// First byte of 128K ROM 1 is also 0xF3 (it's the 48K ROM)
		assert.Equal(t, "ROM 1 byte", bus.Read(0x0000), uint8(0xF3))
	}),
	dsl.NewScenario("Spectrum 128K Paging Lock", func(t *testing.T) {
		bus := NewBus128()
		
		// Lock paging (Port 0x7FFD, bit 5 = 1)
		bus.Out(0x7FFD, 0x20)
		
		// Try to page in RAM 1
		bus.Out(0x7FFD, 0x01)
		
		// Should still be RAM 0
		bus.Write(0xC000, 0xAA)
		bus.Out(0x7FFD, 0x00) // Attempt to page back (should be ignored)
		assert.Equal(t, "RAM should still be locked", bus.Read(0xC000), uint8(0xAA))
	}),
	dsl.NewScenario("Spectrum 128K AY ports", func(t *testing.T) {
		bus := NewBus128()
		
		// Select AY register 7
		bus.Out(0xFFFD, 0x07)
		// Write 0xBF to register 7
		bus.Out(0xBFFD, 0xBF)
		// Read back from register 7
		assert.Equal(t, "AY Register 7", bus.In(0xFFFD), uint8(0xBF))
	}),
	dsl.NewScenario("Spectrum 128K broad port decoding", func(t *testing.T) {
		bus := NewBus128()
		
		// Test 0x7FFD alias (bit 15=0, bit 1=0)
		// 0x00FD should alias to 0x7FFD
		bus.Out(0x00FD, 0x01) // Page RAM 1
		bus.Write(0xC000, 0x55)
		assert.Equal(t, "Value in RAM 1 via alias", bus.Read(0xC000), uint8(0x55))
		
		// Test 0xFFFD alias (bit 15=1, bit 14=1, bit 1=0)
		// 0xC0FD should alias to 0xFFFD
		bus.Out(0xC0FD, 0x01) // Select AY Reg 1
		bus.Out(0x80FD, 0xAA) // 0x80FD alias to 0xBFFD (bit 15=1, bit 14=0, bit 1=0)
		assert.Equal(t, "AY Reg 1 value via alias", bus.In(0xC0FD), uint8(0xAA))
	}),
}
