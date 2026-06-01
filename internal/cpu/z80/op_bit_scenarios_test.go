package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var bitScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: CPL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA // 1010 1010
			bus.Write(0x0000, 0x2F) // CPL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55)) // 0101 0101
			assert.True(t, "Half-Carry should be set", cpu.Regs.Flag(FlagH))
			assert.True(t, "Add/Sub should be set", cpu.Regs.Flag(FlagN))
		},
	},
}
