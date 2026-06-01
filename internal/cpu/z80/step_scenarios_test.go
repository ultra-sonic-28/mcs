package z80

import (
	"mcs/testutils/assert"
	"testing"
)

// stepScenarios defines the test cases for the CPU Step method.
var stepScenarios = []CPUScenario{
	{
		Name: "Step in Halt State",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.SetHalt(true)
			cpu.Regs.PC = 0x1000
			
			cycles := cpu.Step()
			
			assert.Equal(t, "Should return 4 cycles in halt", cycles, 4)
			assert.Equal(t, "Total cycles should be 4", cpu.Cycles, uint64(4))
			assert.Equal(t, "PC should not change in halt", cpu.Regs.PC, uint16(0x1000))
			assert.True(t, "Should still be halted", cpu.Halted)
		},
	},
	{
		Name: "Step Normal Instruction (NOP)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0x00) // NOP
			
			cycles := cpu.Step()
			
			assert.Equal(t, "Should return 4 cycles for NOP", cycles, 4)
			assert.Equal(t, "Total cycles should be 4", cpu.Cycles, uint64(4))
			assert.Equal(t, "PC should be 1", cpu.Regs.PC, uint16(1))
		},
	},
	{
		Name: "Step Instruction with Operand (LD A, n)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0x3E) // LD A, n
			bus.Write(0x0001, 0x42) // value 0x42
			
			cycles := cpu.Step()
			
			assert.Equal(t, "Should return 7 cycles for LD A, n", cycles, 7)
			assert.Equal(t, "Total cycles should be 7", cpu.Cycles, uint64(7))
			assert.Equal(t, "A should be 0x42", cpu.Regs.A, uint8(0x42))
			assert.Equal(t, "PC should be 2", cpu.Regs.PC, uint16(2))
		},
	},
	{
		Name: "Step and Accumulate Cycles",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0x00) // NOP (4 cycles)
			bus.Write(0x0001, 0x3E) // LD A, n (7 cycles)
			bus.Write(0x0002, 0x10)
			
			cpu.Step() // NOP
			assert.Equal(t, "Total cycles after NOP should be 4", cpu.Cycles, uint64(4))
			
			cpu.Step() // LD A, n
			assert.Equal(t, "Total cycles after LD A should be 11", cpu.Cycles, uint64(11))
		},
	},
}
