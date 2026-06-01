package z80

import (
	"mcs/testutils/assert"
	"testing"
)

// fetchScenarios defines the test cases for the CPU FetchByte and FetchWord methods.
var fetchScenarios = []CPUScenario{
	{
		Name: "FetchByte Increments PC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xAB)
			
			val := cpu.FetchByte()
			
			assert.Equal(t, "FetchByte should return correct value", val, uint8(0xAB))
			assert.Equal(t, "PC should be incremented", cpu.Regs.PC, uint16(0x1001))
		},
	},
	{
		Name: "FetchWord (Little-Endian) Increments PC by 2",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x2000
			// Little-Endian: 0x34 0x12 -> 0x1234
			bus.Write(0x2000, 0x34)
			bus.Write(0x2001, 0x12)
			
			val := cpu.FetchWord()
			
			assert.Equal(t, "FetchWord should return correct 16-bit value", val, uint16(0x1234))
			assert.Equal(t, "PC should be incremented by 2", cpu.Regs.PC, uint16(0x2002))
		},
	},
	{
		Name: "Multiple FetchBytes",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.PC = 0x0000
			bus.Write(0x0000, 0x11)
			bus.Write(0x0001, 0x22)
			bus.Write(0x0002, 0x33)
			
			assert.Equal(t, "First FetchByte", cpu.FetchByte(), uint8(0x11))
			assert.Equal(t, "Second FetchByte", cpu.FetchByte(), uint8(0x22))
			assert.Equal(t, "Third FetchByte", cpu.FetchByte(), uint8(0x33))
			assert.Equal(t, "PC after 3 fetches", cpu.Regs.PC, uint16(0x0003))
		},
	},
}
