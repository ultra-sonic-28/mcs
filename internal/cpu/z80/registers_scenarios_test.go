package z80

import (
	"mcs/testutils/assert"
	"testing"
)

// RegisterScenario defines a test case for register operations.
type RegisterScenario struct {
	Name string
	Run  func(t *testing.T)
}

// Scenarios for 16-bit register access.
var register16Scenarios = []RegisterScenario{
	{
		Name: "Set and Get BC",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.SetBC(0x1234)
			assert.Equal(t, "B should be 0x12", r.B, uint8(0x12))
			assert.Equal(t, "C should be 0x34", r.C, uint8(0x34))
			assert.Equal(t, "BC should be 0x1234", r.BC(), uint16(0x1234))
		},
	},
	{
		Name: "Set and Get DE",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.SetDE(0xABCD)
			assert.Equal(t, "D should be 0xAB", r.D, uint8(0xAB))
			assert.Equal(t, "E should be 0xCD", r.E, uint8(0xCD))
			assert.Equal(t, "DE should be 0xABCD", r.DE(), uint16(0xABCD))
		},
	},
	{
		Name: "Set and Get HL",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.SetHL(0x5678)
			assert.Equal(t, "H should be 0x56", r.H, uint8(0x56))
			assert.Equal(t, "L should be 0x78", r.L, uint8(0x78))
			assert.Equal(t, "HL should be 0x5678", r.HL(), uint16(0x5678))
		},
	},
	{
		Name: "Set and Get AF",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.SetAF(0xFA01)
			assert.Equal(t, "A should be 0xFA", r.A, uint8(0xFA))
			assert.Equal(t, "F should be 0x01", r.F, uint8(0x01))
			assert.Equal(t, "AF should be 0xFA01", r.AF(), uint16(0xFA01))
		},
	},
}

// Scenarios for exchange operations.
var exchangeScenarios = []RegisterScenario{
	{
		Name: "Exchange AF and AF'",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.SetAF(0x1234)
			r.SetAFPrime(0xABCD)
			r.ExchangeAF()
			assert.Equal(t, "AF should now be 0xABCD", r.AF(), uint16(0xABCD))
			assert.Equal(t, "AFPrime should now be 0x1234", r.AFPrime(), uint16(0x1234))
		},
	},
	{
		Name: "Exchange Main Registers (EXX)",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.SetBC(0x1111)
			r.SetDE(0x2222)
			r.SetHL(0x3333)
			r.BPrime, r.CPrime = 0xAA, 0xAA
			r.DPrime, r.EPrime = 0xBB, 0xBB
			r.HPrime, r.LPrime = 0xCC, 0xCC

			r.ExchangeMainSwaps()

			assert.Equal(t, "BC should be 0xAAAA", r.BC(), uint16(0xAAAA))
			assert.Equal(t, "DE should be 0xBBBB", r.DE(), uint16(0xBBBB))
			assert.Equal(t, "HL should be 0xCCCC", r.HL(), uint16(0xCCCC))
			assert.Equal(t, "BPrime should be 0x11", r.BPrime, uint8(0x11))
		},
	},
}

// Scenarios for flag operations.
var flagScenarios = []RegisterScenario{
	{
		Name: "Set and Get Flags",
		Run: func(t *testing.T) {
			r := NewRegisters()
			r.SetFlag(FlagZ, true)
			assert.True(t, "Zero flag should be set", r.Flag(FlagZ))
			assert.False(t, "Carry flag should not be set", r.Flag(FlagC))

			r.SetFlag(FlagC, true)
			assert.True(t, "Carry flag should be set", r.Flag(FlagC))

			r.SetFlag(FlagZ, false)
			assert.False(t, "Zero flag should be cleared", r.Flag(FlagZ))
		},
	},
}
