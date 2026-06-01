package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var adcScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: ADC A, n (no Carry)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetFlag(FlagC, false)
			bus.Write(0x0000, 0xCE) // ADC A, n
			bus.Write(0x0001, 0x20)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x30", cpu.Regs.A, uint8(0x30))
			assert.False(t, "Carry should be false", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: ADC A, n (with Carry)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xCE) // ADC A, n
			bus.Write(0x0001, 0x20)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x31", cpu.Regs.A, uint8(0x31))
			assert.False(t, "Carry should be false", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: ADC A, (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xFF
			cpu.Regs.SetFlag(FlagC, true)
			cpu.Regs.SetHL(0x1000)
			bus.Write(0x1000, 0x01)
			bus.Write(0x0000, 0x8E) // ADC A, (HL)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x01", cpu.Regs.A, uint8(0x01))
			assert.True(t, "Carry should be true", cpu.Regs.Flag(FlagC))
			assert.True(t, "Half-Carry should be true", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: ADC A, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.IX = 0x2000
			bus.Write(0x2005, 0x05)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x8E)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x15", cpu.Regs.A, uint8(0x15))
		},
	},
	{
		Name: "Instruction Execution: ADC A, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIXH(0x20)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x8C)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x31", cpu.Regs.A, uint8(0x31))
		},
	},
	{
		Name: "Instruction Execution: ADC A, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIXL(0x05)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x8D)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x16", cpu.Regs.A, uint8(0x16))
		},
	},
	{
		Name: "Instruction Execution: ADC A, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIYH(0x20)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x8C)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x31", cpu.Regs.A, uint8(0x31))
		},
	},
	{
		Name: "Instruction Execution: ADC A, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIYL(0x05)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x8D)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x16", cpu.Regs.A, uint8(0x16))
		},
	},
	{
		Name: "Instruction Execution: ADC HL, BC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x1000)
			cpu.Regs.SetBC(0x2000)
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x4A)
			
			cpu.Step()
			
			assert.Equal(t, "HL should be 0x3001", cpu.Regs.HL(), uint16(0x3001))
			assert.False(t, "Zero flag should be false", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: ADC HL, HL (Zero Result)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x0000)
			cpu.Regs.SetFlag(FlagC, false)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x6A)
			
			cpu.Step()
			
			assert.Equal(t, "HL should be 0x0000", cpu.Regs.HL(), uint16(0x0000))
			assert.True(t, "Zero flag should be true", cpu.Regs.Flag(FlagZ))
		},
	},
}
