package bus

import (
	"mcs/testutils/assert"
	"testing"
)

type BusScenario struct {
	Name string
	Run  func(t *testing.T)
}

var busScenarios = []BusScenario{
	{
		Name: "SimpleBus Memory Read and Write",
		Run: func(t *testing.T) {
			b := NewSimpleBus()
			
			// Test initial state
			assert.Equal(t, "Memory at 0x1234 should be 0 initially", b.Read(0x1234), uint8(0))
			
			// Test write and read back
			b.Write(0x1234, 0x42)
			assert.Equal(t, "Memory read at 0x1234", b.Read(0x1234), uint8(0x42))
			
			// Test boundaries
			b.Write(0x0000, 0x01)
			assert.Equal(t, "Memory at 0x0000", b.Read(0x0000), uint8(0x01))
			
			b.Write(0xFFFF, 0xFF)
			assert.Equal(t, "Memory at 0xFFFF", b.Read(0xFFFF), uint8(0xFF))
		},
	},
	{
		Name: "SimpleBus IO In and Out",
		Run: func(t *testing.T) {
			b := NewSimpleBus()
			
			// Test In (should always return 0 for now)
			assert.Equal(t, "IO In at 0x80", b.In(0x80), uint8(0))
			
			// Test Out (should be a no-op, no panic)
			b.Out(0x80, 0xAA)
			assert.Equal(t, "IO In at 0x80 after Out", b.In(0x80), uint8(0))
		},
	},
}
