package z80

import (
	"mcs/testutils/assert"
	"testing"
)

type AddressingModeScenario struct {
	Name string
	Run  func(t *testing.T)
}

var addressingModeScenarios = []AddressingModeScenario{
	{
		Name: "AddressingMode String Representation",
		Run: func(t *testing.T) {
			cases := []struct {
				mode AddressingMode
				want string
			}{
				{AddrModeNone, "None"},
				{AddrModeImplied, "Implied"},
				{AddrModeAccumulator, "Accumulator"},
				{AddrModeRegister, "Register"},
				{AddrModeRegisterPair, "Register Pair"},
				{AddrModeImmediate, "Immediate"},
				{AddrModeImmediate16, "Immediate 16"},
				{AddrModeIndirect, "Indirect"},
				{AddrModeExtended, "Extended"},
				{AddrModeIndexed, "Indexed"},
				{AddrModeRelative, "Relative"},
				{AddrModePort, "Port"},
				{AddrModeBit, "Bit"},
				{AddressingMode(99), "Unknown"},
			}

			for _, c := range cases {
				assert.Equal(t, "Mode string should match", c.mode.String(), c.want)
			}
		},
	},
}
