package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var ioScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: OUT (n), A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0x12
			bus.Write(0x0000, 0xD3) // OUT (n), A
			bus.Write(0x0001, 0x34) // n = 0x34
			
			cpu.Step()
			
			assert.Equal(t, "IO port 0x1234 should have 0x12", bus.IO[0x1234], uint8(0x12))
			assert.Equal(t, "Last IO address should be 0x1234", bus.Last, uint16(0x1234))
		},
	},
	{
		Name: "Instruction Execution: IN A, (n)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0xAB
			bus.IO[0xAB42] = 0x55
			bus.Write(0x0000, 0xDB) // IN A, (n)
			bus.Write(0x0001, 0x42) // n = 0x42
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
			assert.Equal(t, "Last IO address should be 0xAB42", bus.Last, uint16(0xAB42))
		},
	},
	{
		Name: "Instruction Execution: IN B, (C)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.SetBC(0x1234)
			bus.IO[0x1234] = 0x55
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x40) // IN B, (C)
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0x55", cpu.Regs.B, uint8(0x55))
			assert.Equal(t, "Last IO address should be 0x1234", bus.Last, uint16(0x1234))
			assert.False(t, "Sign flag should be cleared", cpu.Regs.Flag(FlagS))
			assert.False(t, "Zero flag should be cleared", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: IN (C) - Special case",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.SetBC(0x0000)
			bus.IO[0x0000] = 0x00
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x70) // IN (C)
			
			cpu.Step()
			
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: OUT (C), D",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.SetBC(0x5678)
			cpu.Regs.D = 0xAA
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x51) // OUT (C), D
			
			cpu.Step()
			
			assert.Equal(t, "IO port 0x5678 should have 0xAA", bus.IO[0x5678], uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: OUT (C), 0 - Special case",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.SetBC(0x0000)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x71) // OUT (C), 0
			
			cpu.Step()
			
			assert.Equal(t, "IO port 0x0000 should have 0", bus.IO[0x0000], uint8(0))
		},
	},
}
