package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var rotScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: RLCA",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x81 // 1000 0001
			bus.Write(0x0000, 0x07) // RLCA
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x03", cpu.Regs.A, uint8(0x03)) // 0000 0011
			assert.True(t, "Carry should be set", cpu.Regs.Flag(FlagC))
			assert.False(t, "Half-Carry should be cleared", cpu.Regs.Flag(FlagH))
			assert.False(t, "Add/Sub should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: RRCA",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x81 // 1000 0001
			bus.Write(0x0000, 0x0F) // RRCA
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xC0", cpu.Regs.A, uint8(0xC0)) // 1100 0000
			assert.True(t, "Carry should be set", cpu.Regs.Flag(FlagC))
			assert.False(t, "Half-Carry should be cleared", cpu.Regs.Flag(FlagH))
			assert.False(t, "Add/Sub should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: RLA",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x81 // 1000 0001
			cpu.Regs.SetFlag(FlagC, false)
			bus.Write(0x0000, 0x17) // RLA
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x02", cpu.Regs.A, uint8(0x02)) // 0000 0010 (bit 0 is old carry)
			assert.True(t, "Carry should be set", cpu.Regs.Flag(FlagC))
			
			// Second test with carry set
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0x17)
			cpu.Regs.A = 0x01
			cpu.Regs.SetFlag(FlagC, true)
			
			cpu.Step()
			assert.Equal(t, "A should be 0x03", cpu.Regs.A, uint8(0x03)) // 0000 0011 (bit 0 is old carry)
			assert.False(t, "Carry should be cleared", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: RRA",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x81 // 1000 0001
			cpu.Regs.SetFlag(FlagC, false)
			bus.Write(0x0000, 0x1F) // RRA
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x40", cpu.Regs.A, uint8(0x40)) // 0100 0000 (bit 7 is old carry)
			assert.True(t, "Carry should be set", cpu.Regs.Flag(FlagC))
			
			// Second test with carry set
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0x1F)
			cpu.Regs.A = 0x80
			cpu.Regs.SetFlag(FlagC, true)
			
			cpu.Step()
			assert.Equal(t, "A should be 0xC0", cpu.Regs.A, uint8(0xC0)) // 1100 0000 (bit 7 is old carry)
			assert.False(t, "Carry should be cleared", cpu.Regs.Flag(FlagC))
		},
	},
}
