package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var miscScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: SCF",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Initial state: C=0, H=1, N=1
			cpu.Regs.SetFlag(FlagC, false)
			cpu.Regs.SetFlag(FlagH, true)
			cpu.Regs.SetFlag(FlagN, true)
			
			bus.Write(0x0000, 0x37) // SCF
			cpu.Step()
			
			assert.True(t, "Carry flag should be set", cpu.Regs.Flag(FlagC))
			assert.False(t, "Half-Carry flag should be cleared", cpu.Regs.Flag(FlagH))
			assert.False(t, "Add/Sub flag should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: CCF (Carry was 0)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Initial state: C=0, H=1, N=1
			cpu.Regs.SetFlag(FlagC, false)
			cpu.Regs.SetFlag(FlagH, true)
			cpu.Regs.SetFlag(FlagN, true)
			
			bus.Write(0x0000, 0x3F) // CCF
			cpu.Step()
			
			assert.True(t, "Carry flag should be complemented (0->1)", cpu.Regs.Flag(FlagC))
			assert.False(t, "Half-Carry flag should take old Carry state (0)", cpu.Regs.Flag(FlagH))
			assert.False(t, "Add/Sub flag should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: CCF (Carry was 1)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Initial state: C=1, H=0, N=1
			cpu.Regs.SetFlag(FlagC, true)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, true)
			
			bus.Write(0x0000, 0x3F) // CCF
			cpu.Step()
			
			assert.False(t, "Carry flag should be complemented (1->0)", cpu.Regs.Flag(FlagC))
			assert.True(t, "Half-Carry flag should take old Carry state (1)", cpu.Regs.Flag(FlagH))
			assert.False(t, "Add/Sub flag should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: 0xDD (NOP in MainTable)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			instr := MainTable[0xDD]
			cycles := instr.Execute(cpu)
			
			assert.Equal(t, "Mnemonic should be NOP", instr.Mnemonic, "NOP")
			assert.Equal(t, "Cycles should be 4", cycles, 4)
		},
	},
	{
		Name: "Instruction Execution: 0xFD (NOP in MainTable)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			instr := MainTable[0xFD]
			cycles := instr.Execute(cpu)
			
			assert.Equal(t, "Mnemonic should be NOP", instr.Mnemonic, "NOP")
			assert.Equal(t, "Cycles should be 4", cycles, 4)
		},
	},
	{
		Name: "Instruction Execution: DI",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.IFF1 = true
			cpu.IFF2 = true
			bus.Write(0x0000, 0xF3) // DI
			
			cpu.Step()
			
			assert.False(t, "IFF1 should be false", cpu.IFF1)
			assert.False(t, "IFF2 should be false", cpu.IFF2)
		},
	},
	{
		Name: "Instruction Execution: EI",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.IFF1 = false
			cpu.IFF2 = false
			bus.Write(0x0000, 0xFB) // EI
			
			cpu.Step()
			
			assert.True(t, "IFF1 should be true", cpu.IFF1)
			assert.True(t, "IFF2 should be true", cpu.IFF2)
		},
	},
}
