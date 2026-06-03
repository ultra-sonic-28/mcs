// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

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
		// It might be slightly more due to the last instruction overshoot.
		assert.True(t, "Total cycles should be >= 69888", m.CPU.Cycles >= 69888)
		assert.True(t, "Total cycles should be close to 69888", m.CPU.Cycles < 69900)
	}),
}
