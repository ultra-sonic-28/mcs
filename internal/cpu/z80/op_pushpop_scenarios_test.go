package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var pushPopScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: PUSH BC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0x1000
			cpu.Regs.SetBC(0x1234)
			bus.Write(0x0000, 0xC5) // PUSH BC
			
			cpu.Step()
			
			assert.Equal(t, "SP should be 0x0FFE", cpu.Regs.SP, uint16(0x0FFE))
			assert.Equal(t, "Mem[0x0FFE] should be 0x34", bus.Read(0x0FFE), uint8(0x34))
			assert.Equal(t, "Mem[0x0FFF] should be 0x12", bus.Read(0x0FFF), uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: POP BC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0x0FFE
			bus.Write(0x0FFE, 0x34)
			bus.Write(0x0FFF, 0x12)
			bus.Write(0x0000, 0xC1) // POP BC
			
			cpu.Step()
			
			assert.Equal(t, "BC should be 0x1234", cpu.Regs.BC(), uint16(0x1234))
			assert.Equal(t, "SP should be 0x1000", cpu.Regs.SP, uint16(0x1000))
		},
	},
	{
		Name: "Instruction Execution: PUSH AF",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0x1000
			cpu.Regs.SetAF(0x55AA)
			bus.Write(0x0000, 0xF5) // PUSH AF
			
			cpu.Step()
			
			assert.Equal(t, "SP should be 0x0FFE", cpu.Regs.SP, uint16(0x0FFE))
			assert.Equal(t, "Mem[0x0FFE] should be 0xAA", bus.Read(0x0FFE), uint8(0xAA))
			assert.Equal(t, "Mem[0x0FFF] should be 0x55", bus.Read(0x0FFF), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: POP AF (Flags)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0x0FFE
			bus.Write(0x0FFE, 0xFF) // Flags
			bus.Write(0x0FFF, 0x55) // Accumulator
			bus.Write(0x0000, 0xF1) // POP AF
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
			assert.Equal(t, "F should be 0xFF", cpu.Regs.F, uint8(0xFF))
			assert.True(t, "Carry should be true", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: PUSH IX",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0x1000
			cpu.Regs.IX = 0xABCD
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xE5) // PUSH IX
			
			cpu.Step()
			
			assert.Equal(t, "SP should be 0x0FFE", cpu.Regs.SP, uint16(0x0FFE))
			assert.Equal(t, "Mem[0x0FFE] should be 0xCD", bus.Read(0x0FFE), uint8(0xCD))
			assert.Equal(t, "Mem[0x0FFF] should be 0xAB", bus.Read(0x0FFF), uint8(0xAB))
		},
	},
	{
		Name: "Instruction Execution: POP IY",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0x0FFE
			bus.Write(0x0FFE, 0x78)
			bus.Write(0x0FFF, 0x56)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xE1) // POP IY
			
			cpu.Step()
			
			assert.Equal(t, "IY should be 0x5678", cpu.Regs.IY, uint16(0x5678))
			assert.Equal(t, "SP should be 0x1000", cpu.Regs.SP, uint16(0x1000))
		},
	},
}
