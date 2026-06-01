package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var jumpScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: JP nn",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0xC3) // JP nn
			bus.Write(0x0001, 0x34) // low
			bus.Write(0x0002, 0x12) // high
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x1234", cpu.Regs.PC, uint16(0x1234))
		},
	},
	{
		Name: "Instruction Execution: JP Z, nn (Jump Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetFlag(FlagZ, true)
			bus.Write(0x0000, 0xCA) // JP Z, nn
			bus.Write(0x0001, 0x00)
			bus.Write(0x0002, 0x20)
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x2000", cpu.Regs.PC, uint16(0x2000))
		},
	},
	{
		Name: "Instruction Execution: JP Z, nn (Jump Not Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetFlag(FlagZ, false)
			bus.Write(0x0000, 0xCA) // JP Z, nn
			bus.Write(0x0001, 0x00)
			bus.Write(0x0002, 0x20)
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 3", cpu.Regs.PC, uint16(3))
		},
	},
	{
		Name: "Instruction Execution: JR e",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0x18) // JR e
			bus.Write(0x0001, 0x05) // e = 5
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 7", cpu.Regs.PC, uint16(7))
		},
	},
	{
		Name: "Instruction Execution: CALL nn and RET",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0xFFFF
			
			// Main: CALL 0x1000
			bus.Write(0x0000, 0xCD) // CALL nn
			bus.Write(0x0001, 0x00)
			bus.Write(0x0002, 0x10)
			
			// Subroutine: RET
			bus.Write(0x1000, 0xC9) // RET
			
			cpu.Step() // Execute CALL
			
			assert.Equal(t, "PC should be 0x1000", cpu.Regs.PC, uint16(0x1000))
			assert.Equal(t, "SP should be 0xFFFD", cpu.Regs.SP, uint16(0xFFFD))
			assert.Equal(t, "Stacked PC Low should be 0x03", bus.Read(0xFFFD), uint8(0x03))
			assert.Equal(t, "Stacked PC High should be 0x00", bus.Read(0xFFFE), uint8(0x00))
			
			cpu.Step() // Execute RET
			
			assert.Equal(t, "PC should be 0x0003", cpu.Regs.PC, uint16(0x0003))
			assert.Equal(t, "SP should be 0xFFFF", cpu.Regs.SP, uint16(0xFFFF))
		},
	},
	{
		Name: "Instruction Execution: JR NZ, e (Jump Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetFlag(FlagZ, false)
			bus.Write(0x0000, 0x20) // JR NZ, e
			bus.Write(0x0001, 0x05) // e = 5
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 7", cpu.Regs.PC, uint16(7))
		},
	},
	{
		Name: "Instruction Execution: JR NZ, e (Jump Not Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetFlag(FlagZ, true)
			bus.Write(0x0000, 0x20) // JR NZ, e
			bus.Write(0x0001, 0x05)
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 2", cpu.Regs.PC, uint16(2))
		},
	},
	{
		Name: "Instruction Execution: CALL C, nn (Call Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetFlag(FlagC, true)
			cpu.Regs.SP = 0xFFFF
			bus.Write(0x0000, 0xDC) // CALL C, nn
			bus.Write(0x0001, 0x00)
			bus.Write(0x0002, 0x10)
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x1000", cpu.Regs.PC, uint16(0x1000))
			assert.Equal(t, "SP should be 0xFFFD", cpu.Regs.SP, uint16(0xFFFD))
		},
	},
	{
		Name: "Instruction Execution: RET Z (Return Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetFlag(FlagZ, true)
			cpu.Regs.SP = 0xFFFD
			bus.Write(0xFFFD, 0x34)
			bus.Write(0xFFFE, 0x12)
			bus.Write(0x0000, 0xC8) // RET Z
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x1234", cpu.Regs.PC, uint16(0x1234))
			assert.Equal(t, "SP should be 0xFFFF", cpu.Regs.SP, uint16(0xFFFF))
		},
	},
	{
		Name: "Instruction Execution: DJNZ e (Jump Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.B = 0x05
			bus.Write(0x0000, 0x10) // DJNZ e
			bus.Write(0x0001, 0x0A) // e = 10
			
			cycles := cpu.Step()
			
			assert.Equal(t, "B should be 4", cpu.Regs.B, uint8(4))
			assert.Equal(t, "PC should be 12", cpu.Regs.PC, uint16(12))
			assert.Equal(t, "Cycles should be 13", cycles, 13)
		},
	},
	{
		Name: "Instruction Execution: DJNZ e (Jump Not Taken)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.B = 0x01
			bus.Write(0x0000, 0x10) // DJNZ e
			bus.Write(0x0001, 0x0A)
			
			cycles := cpu.Step()
			
			assert.Equal(t, "B should be 0", cpu.Regs.B, uint8(0))
			assert.Equal(t, "PC should be 2", cpu.Regs.PC, uint16(2))
			assert.Equal(t, "Cycles should be 8", cycles, 8)
		},
	},
	{
		Name: "Instruction Execution: RST 00H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1234
			cpu.Regs.SP = 0xFFFF
			bus.Write(0x1234, 0xC7) // RST 00H
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x0000", cpu.Regs.PC, uint16(0x0000))
			assert.Equal(t, "SP should be 0xFFFD", cpu.Regs.SP, uint16(0xFFFD))
			assert.Equal(t, "Stacked PC Low should be 0x35", bus.Read(0xFFFD), uint8(0x35)) // PC was 0x1234, incremented to 0x1235 before jump? Wait, Step() increments after fetch.
			// FetchByte increments PC. So after fetching 0xC7, PC is 0x1235.
		},
	},
	{
		Name: "Instruction Execution: RST 08H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x2000
			cpu.Regs.SP = 0x1000
			bus.Write(0x2000, 0xCF) // RST 08H
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x0008", cpu.Regs.PC, uint16(0x0008))
			assert.Equal(t, "SP should be 0x0FFE", cpu.Regs.SP, uint16(0x0FFE))
		},
	},
	{
		Name: "Instruction Execution: RST 10H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xD7) // RST 10H
			cpu.Step()
			assert.Equal(t, "PC should be 0x0010", cpu.Regs.PC, uint16(0x0010))
		},
	},
	{
		Name: "Instruction Execution: RST 18H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xDF) // RST 18H
			cpu.Step()
			assert.Equal(t, "PC should be 0x0018", cpu.Regs.PC, uint16(0x0018))
		},
	},
	{
		Name: "Instruction Execution: RST 20H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xE7) // RST 20H
			cpu.Step()
			assert.Equal(t, "PC should be 0x0020", cpu.Regs.PC, uint16(0x0020))
		},
	},
	{
		Name: "Instruction Execution: RST 28H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xEF) // RST 28H
			cpu.Step()
			assert.Equal(t, "PC should be 0x0028", cpu.Regs.PC, uint16(0x0028))
		},
	},
	{
		Name: "Instruction Execution: RST 30H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xF7) // RST 30H
			cpu.Step()
			assert.Equal(t, "PC should be 0x0030", cpu.Regs.PC, uint16(0x0030))
		},
	},
	{
		Name: "Instruction Execution: RST 38H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xFF) // RST 38H
			cpu.Step()
			assert.Equal(t, "PC should be 0x0038", cpu.Regs.PC, uint16(0x0038))
		},
	},
	{
		Name: "Instruction Execution: JP (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x1234)
			bus.Write(0x0000, 0xE9) // JP (HL)
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x1234", cpu.Regs.PC, uint16(0x1234))
		},
	},
	{
		Name: "Instruction Execution: JP (IX)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x5678
			bus.Write(0x0000, 0xDD) // IX prefix
			bus.Write(0x0001, 0xE9) // JP (IX)
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0x5678", cpu.Regs.PC, uint16(0x5678))
		},
	},
	{
		Name: "Instruction Execution: JP (IY)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0xABCD
			bus.Write(0x0000, 0xFD) // IY prefix
			bus.Write(0x0001, 0xE9) // JP (IY)
			
			cpu.Step()
			
			assert.Equal(t, "PC should be 0xABCD", cpu.Regs.PC, uint16(0xABCD))
		},
	},
}
