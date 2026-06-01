package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var exchangeInstructionScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: EX DE, HL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetDE(0x1234)
			cpu.Regs.SetHL(0xABCD)
			bus.Write(0x0000, 0xEB) // EX DE, HL
			
			cpu.Step()
			
			assert.Equal(t, "DE should be 0xABCD", cpu.Regs.DE(), uint16(0xABCD))
			assert.Equal(t, "HL should be 0x1234", cpu.Regs.HL(), uint16(0x1234))
		},
	},
	{
		Name: "Instruction Execution: EX AF, AF'",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetAF(0x1122)
			cpu.Regs.SetAFPrime(0x3344)
			bus.Write(0x0000, 0x08) // EX AF, AF'
			
			cpu.Step()
			
			assert.Equal(t, "AF should be 0x3344", cpu.Regs.AF(), uint16(0x3344))
			assert.Equal(t, "AF' should be 0x1122", cpu.Regs.AFPrime(), uint16(0x1122))
		},
	},
	{
		Name: "Instruction Execution: EXX",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetBC(0x1111)
			cpu.Regs.SetDE(0x2222)
			cpu.Regs.SetHL(0x3333)
			cpu.Regs.BPrime, cpu.Regs.CPrime = 0xAA, 0xBB
			cpu.Regs.DPrime, cpu.Regs.EPrime = 0xCC, 0xDD
			cpu.Regs.HPrime, cpu.Regs.LPrime = 0xEE, 0xFF
			bus.Write(0x0000, 0xD9) // EXX
			
			cpu.Step()
			
			assert.Equal(t, "BC should be swapped", cpu.Regs.BC(), uint16(0xAABB))
			assert.Equal(t, "DE should be swapped", cpu.Regs.DE(), uint16(0xCCDD))
			assert.Equal(t, "HL should be swapped", cpu.Regs.HL(), uint16(0xEEFF))
		},
	},
	{
		Name: "Instruction Execution: EX (SP), HL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SP = 0x1000
			cpu.Regs.SetHL(0xABCD)
			bus.Write(0x1000, 0x34)
			bus.Write(0x1001, 0x12)
			bus.Write(0x0000, 0xE3) // EX (SP), HL
			
			cpu.Step()
			
			assert.Equal(t, "HL should now be 0x1234", cpu.Regs.HL(), uint16(0x1234))
			assert.Equal(t, "Memory at SP should be CD", bus.Read(0x1000), uint8(0xCD))
			assert.Equal(t, "Memory at SP+1 should be AB", bus.Read(0x1001), uint8(0xAB))
		},
	},
}
