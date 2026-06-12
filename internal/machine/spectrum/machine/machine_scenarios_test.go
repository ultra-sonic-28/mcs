// Package machine implements the ZX Spectrum machine logic.
package machine

import (
	"mcs/internal/cpu/z80"
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var machineScenarios = []dsl.Scenario{
	dsl.NewScenario("Machine RunFrame execution", func(t *testing.T) {
		m := NewMachine()
		m.Reset()

		// Set PC to a HALT instruction to avoid executing random memory
		m.Bus.Write(0x0000, 0x76) // HALT
		m.CPU.IFF1 = true
		m.CPU.IM = z80.IM1

		m.RunFrame()

		// T-cycles should be at least CyclesPerFrame (69888)
		assert.True(t, "Total cycles should be >= 69888", m.CPU.Cycles >= 69888)
		assert.True(t, "Total cycles should be close to 69888", m.CPU.Cycles < 69900)
	}),
	dsl.NewScenario("Machine128 RunFrame execution", func(t *testing.T) {
		m := NewMachine128()
		m.Reset()

		m.Bus.Write(0x0000, 0x76) // HALT
		m.CPU.IFF1 = true
		m.CPU.IM = z80.IM1

		m.RunFrame()

		// T-cycles should be at least CyclesPerFrame128 (70908)
		assert.True(t, "Total cycles should be >= 70908", m.CPU.Cycles >= 70908)
		assert.True(t, "Total cycles should be close to 70908", m.CPU.Cycles < 70920)
	}),

}
