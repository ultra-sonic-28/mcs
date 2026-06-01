package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var subcScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: SBC A, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIXH(0x10)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x9C)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x1F", cpu.Regs.A, uint8(0x1F)) // 0x30 - 0x10 - 1
		},
	},
	{
		Name: "Instruction Execution: SBC A, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.SetIXL(0x05)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x9D)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x0F", cpu.Regs.A, uint8(0x0F)) // 0x15 - 0x05 - 1
		},
	},
	{
		Name: "Instruction Execution: SBC A, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.IX = 0x2000
			bus.Write(0x2005, 0x05)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x9E)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x0F", cpu.Regs.A, uint8(0x0F)) // 0x15 - 0x05 - 1
		},
	},
	{
		Name: "Instruction Execution: SBC A, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIYH(0x10)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x9C)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x1F", cpu.Regs.A, uint8(0x1F)) // 0x30 - 0x10 - 1
		},
	},
	{
		Name: "Instruction Execution: SBC A, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.SetIYL(0x05)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x9D)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x0F", cpu.Regs.A, uint8(0x0F)) // 0x15 - 0x05 - 1
		},
	},
	{
		Name: "Instruction Execution: SBC A, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.IY = 0x2000
			bus.Write(0x2005, 0x05)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x9E)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x0F", cpu.Regs.A, uint8(0x0F)) // 0x15 - 0x05 - 1
		},
	},
	{
		Name: "Instruction Execution: SBC HL, BC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x1000)
			cpu.Regs.SetBC(0x0100)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x42) // SBC HL, BC
			
			cpu.Step()
			
			assert.Equal(t, "HL should be 0x0EFF", cpu.Regs.HL(), uint16(0x0EFF)) // 0x1000 - 0x0100 - 1
			assert.True(t, "Add/Subtract flag should be set", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: SBC HL, DE",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x0001)
			cpu.Regs.SetDE(0x0001)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x52) // SBC HL, DE
			
			cpu.Step()
			
			assert.Equal(t, "HL should be 0xFFFF", cpu.Regs.HL(), uint16(0xFFFF)) // 0x0001 - 0x0001 - 1
			assert.True(t, "Carry flag should be set (borrow)", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: SBC HL, HL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x1234)
			cpu.Regs.SetFlag(FlagC, false)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x62) // SBC HL, HL
			
			cpu.Step()
			
			assert.Equal(t, "HL should be 0x0000", cpu.Regs.HL(), uint16(0x0000))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: SBC HL, SP",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x2000)
			cpu.Regs.SP = 0x1000
			cpu.Regs.SetFlag(FlagC, false)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x72) // SBC HL, SP
			
			cpu.Step()
			
			assert.Equal(t, "HL should be 0x1000", cpu.Regs.HL(), uint16(0x1000))
		},
	},
}
