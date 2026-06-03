package z80

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var interruptScenarios = []dsl.Scenario{
	dsl.NewScenario("IM1 Interrupt handling", func(t *testing.T) {
		bus := &MockBus{}
		cpu := NewCPU(bus, bus)
		cpu.Regs.PC = 0x1000
		cpu.Regs.SP = 0x2000
		cpu.IM = IM1
		cpu.IFF1 = true
		cpu.INT = true

		cycles := cpu.HandleInterrupts()
		assert.Equal(t, "IM1 interrupt cycles", cycles, 13)
		assert.Equal(t, "PC after interrupt", cpu.Regs.PC, uint16(0x0038))
		assert.Equal(t, "SP after interrupt", cpu.Regs.SP, uint16(0x1FFE))
		assert.False(t, "IFF1 after interrupt", cpu.IFF1)
		
		// Check pushed PC
		low := bus.Read(0x1FFE)
		high := bus.Read(0x1FFF)
		assert.Equal(t, "Pushed PC", uint16(high)<<8|uint16(low), uint16(0x1000))
	}),
	dsl.NewScenario("NMI handling", func(t *testing.T) {
		bus := &MockBus{}
		cpu := NewCPU(bus, bus)
		cpu.Regs.PC = 0x1000
		cpu.Regs.SP = 0x2000
		cpu.IFF1 = true
		cpu.NMI = true

		cycles := cpu.HandleInterrupts()
		assert.Equal(t, "NMI cycles", cycles, 11)
		assert.Equal(t, "PC after NMI", cpu.Regs.PC, uint16(0x0066))
		assert.False(t, "IFF1 after NMI", cpu.IFF1)
		assert.True(t, "IFF2 should store old IFF1", cpu.IFF2)
	}),
	dsl.NewScenario("Interrupt in Halt state", func(t *testing.T) {
		bus := &MockBus{}
		cpu := NewCPU(bus, bus)
		cpu.SetHalt(true)
		cpu.IM = IM1
		cpu.IFF1 = true
		cpu.INT = true

		cycles := cpu.Step()
		assert.Equal(t, "Interrupt from halt cycles", cycles, 13)
		assert.False(t, "Should not be halted anymore", cpu.Halted)
		assert.Equal(t, "PC after interrupt", cpu.Regs.PC, uint16(0x0038))
	}),
}
