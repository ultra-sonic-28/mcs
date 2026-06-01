package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var logicScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: AND r",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0xAA
			cpu.Regs.B = 0x0F
			bus.Write(0x0000, 0xA0) // AND B
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x0A", cpu.Regs.A, uint8(0x0A))
			assert.True(t, "Half-Carry should be set for AND", cpu.Regs.Flag(FlagH))
			assert.False(t, "Carry should be cleared", cpu.Regs.Flag(FlagC))
			assert.False(t, "Add/Sub should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: AND n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0x55
			bus.Write(0x0000, 0xE6) // AND n
			bus.Write(0x0001, 0x11)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x11", cpu.Regs.A, uint8(0x11))
		},
	},
	{
		Name: "Instruction Execution: AND (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0xFF
			cpu.Regs.SetHL(0x1000)
			bus.Write(0x1000, 0x80)
			bus.Write(0x0000, 0xA6) // AND (HL)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x80", cpu.Regs.A, uint8(0x80))
		},
	},
	{
		Name: "Instruction Execution: OR r",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0xF0
			cpu.Regs.B = 0x0F
			bus.Write(0x0000, 0xB0) // OR B
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
			assert.False(t, "Half-Carry should be cleared for OR", cpu.Regs.Flag(FlagH))
			assert.False(t, "Carry should be cleared", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: OR n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0x11
			bus.Write(0x0000, 0xF6) // OR n
			bus.Write(0x0001, 0x22)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x33", cpu.Regs.A, uint8(0x33))
		},
	},
	{
		Name: "Instruction Execution: OR (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0x00
			cpu.Regs.SetHL(0x2000)
			bus.Write(0x2000, 0x5A)
			bus.Write(0x0000, 0xB6) // OR (HL)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x5A", cpu.Regs.A, uint8(0x5A))
		},
	},
	{
		Name: "Instruction Execution: XOR r",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0x55
			cpu.Regs.B = 0x55
			bus.Write(0x0000, 0xA8) // XOR B
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
			assert.False(t, "Half-Carry should be cleared for XOR", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: XOR n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0xFF
			bus.Write(0x0000, 0xEE) // XOR n
			bus.Write(0x0001, 0x0F)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xF0", cpu.Regs.A, uint8(0xF0))
		},
	},
	{
		Name: "Instruction Execution: XOR (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.A = 0xAA
			cpu.Regs.SetHL(0x3000)
			bus.Write(0x3000, 0xFF)
			bus.Write(0x0000, 0xAE) // XOR (HL)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: AND IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			cpu.Regs.SetIXH(0xF0)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xA4) // AND IXH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xA0", cpu.Regs.A, uint8(0xA0))
			assert.True(t, "Half-Carry should be set for AND", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: AND IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			cpu.Regs.SetIXL(0x0F)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xA5) // AND IXL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x05", cpu.Regs.A, uint8(0x05))
		},
	},
	{
		Name: "Instruction Execution: AND (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xFF
			cpu.Regs.IX = 0x4000
			bus.Write(0x4005, 0x55)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xA6) // AND (IX+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: XOR IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			cpu.Regs.SetIXH(0x55)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xAC) // XOR IXH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
			assert.False(t, "Half-Carry should be cleared for XOR", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: XOR IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			cpu.Regs.SetIXL(0x55)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xAD) // XOR IXL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: XOR (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			cpu.Regs.IX = 0x5000
			bus.Write(0x5005, 0xAA)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xAE) // XOR (IX+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
		},
	},
	{
		Name: "Instruction Execution: OR IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIXH(0x01)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xB4) // OR IXH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x11", cpu.Regs.A, uint8(0x11))
			assert.False(t, "Half-Carry should be cleared for OR", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: OR IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xF0
			cpu.Regs.SetIXL(0x0F)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xB5) // OR IXL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
		},
	},
	{
		Name: "Instruction Execution: OR (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			cpu.Regs.IX = 0x6000
			bus.Write(0x600A, 0xAA)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xB6) // OR (IX+d)
			bus.Write(0x0002, 0x0A) // d = 10
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
		},
	},
	{
		Name: "Instruction Execution: AND IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			cpu.Regs.SetIYH(0xF0)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xA4) // AND IYH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xA0", cpu.Regs.A, uint8(0xA0))
			assert.True(t, "Half-Carry should be set for AND", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: AND IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			cpu.Regs.SetIYL(0x0F)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xA5) // AND IYL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x05", cpu.Regs.A, uint8(0x05))
		},
	},
	{
		Name: "Instruction Execution: AND (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xFF
			cpu.Regs.IY = 0x4000
			bus.Write(0x4005, 0x55)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xA6) // AND (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: XOR IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			cpu.Regs.SetIYH(0x55)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xAC) // XOR IYH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
		},
	},
	{
		Name: "Instruction Execution: XOR IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			cpu.Regs.SetIYL(0x55)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xAD) // XOR IYL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: XOR (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			cpu.Regs.IY = 0x5000
			bus.Write(0x5005, 0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xAE) // XOR (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
		},
	},
	{
		Name: "Instruction Execution: OR IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIYH(0x01)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xB4) // OR IYH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x11", cpu.Regs.A, uint8(0x11))
		},
	},
	{
		Name: "Instruction Execution: OR IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xF0
			cpu.Regs.SetIYL(0x0F)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xB5) // OR IYL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
		},
	},
	{
		Name: "Instruction Execution: OR (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			cpu.Regs.IY = 0x6000
			bus.Write(0x600A, 0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xB6) // OR (IY+d)
			bus.Write(0x0002, 0x0A) // d = 10
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
		},
	},
}
