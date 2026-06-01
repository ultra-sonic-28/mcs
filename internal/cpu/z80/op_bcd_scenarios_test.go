package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var bcdScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: DAA (Addition 0x09+0x01)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			// Mock result of 0x09 + 0x01
			cpu.Regs.A = 0x0A
			cpu.Regs.SetFlag(FlagH, true)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(FlagC, false)
			
			bus.Write(0x0000, 0x27) // DAA
			cpu.Step()
			
			assert.Equal(t, "A should be 0x10", cpu.Regs.A, uint8(0x10))
			assert.False(t, "Carry should be false", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: DAA (Addition 0x99+0x01)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			// Mock result of 0x99 + 0x01 = 0x9A (H=1)
			cpu.Regs.A = 0x9A
			cpu.Regs.SetFlag(FlagH, true)
			cpu.Regs.SetFlag(FlagN, false)
			
			bus.Write(0x0000, 0x27) // DAA
			cpu.Step()
			
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
			assert.True(t, "Carry should be true", cpu.Regs.Flag(FlagC))
			assert.True(t, "Zero flag should be true", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: DAA (Subtraction 0x10-0x01)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			// Mock result of 0x10 - 0x01 = 0x0F (H=1, N=1)
			cpu.Regs.A = 0x0F
			cpu.Regs.SetFlag(FlagH, true)
			cpu.Regs.SetFlag(FlagN, true)
			
			bus.Write(0x0000, 0x27) // DAA
			cpu.Step()
			
			assert.Equal(t, "A should be 0x09", cpu.Regs.A, uint8(0x09))
		},
	},
	{
		Name: "Instruction Execution: RLD",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x5A
			cpu.Regs.SetHL(0x2000)
			bus.Write(0x2000, 0x12)
			
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x6F) // RLD
			
			cpu.Step()
			
			// A[3:0] was 0xA, (HL) was 0x12
			// New (HL) = (HL)[3:0] << 4 | A[3:0] = 0x2A
			// New A[3:0] = (HL)[7:4] = 0x1
			assert.Equal(t, "A should be 0x51", cpu.Regs.A, uint8(0x51))
			assert.Equal(t, "Memory at 0x2000 should be 0x2A", bus.Read(0x2000), uint8(0x2A))
		},
	},
	{
		Name: "Instruction Execution: RRD",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x5A
			cpu.Regs.SetHL(0x3000)
			bus.Write(0x3000, 0x12)
			
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x67) // RRD
			
			cpu.Step()
			
			// A[3:0] was 0xA, (HL) was 0x12
			// New (HL) = A[3:0] << 4 | (HL)[7:4] = 0xA1
			// New A[3:0] = (HL)[3:0] = 0x2
			assert.Equal(t, "A should be 0x52", cpu.Regs.A, uint8(0x52))
			assert.Equal(t, "Memory at 0x3000 should be 0xA1", bus.Read(0x3000), uint8(0xA1))
		},
	},
}
