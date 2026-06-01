package z80

import (
	"mcs/testutils/assert"
	"testing"
)

type CPUScenario struct {
	Name string
	Run  func(t *testing.T)
}

type MockBus struct {
	Mem  [65536]uint8
	IO   [65536]uint8
	Last uint16
}

func (m *MockBus) Read(addr uint16) uint8        { return m.Mem[addr] }
func (m *MockBus) Write(addr uint16, val uint8) { m.Mem[addr] = val }
func (m *MockBus) In(port uint16) uint8         { m.Last = port; return m.IO[port] }
func (m *MockBus) Out(port uint16, val uint8)   { m.IO[port] = val; m.Last = port }

var cpuScenarios = []CPUScenario{
	{
		Name: "CPU Initialization and Reset",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Check default reset state
			assert.Equal(t, "PC should be 0", cpu.Regs.PC, uint16(0))
			assert.Equal(t, "SP should be 0xFFFF", cpu.Regs.SP, uint16(0xFFFF))
			assert.False(t, "IFF1 should be false", cpu.IFF1)
			assert.Equal(t, "IM should be IM0", cpu.IM, IM0)
			assert.Equal(t, "Cycles should be 0", cpu.Cycles, uint64(0))
		},
	},
	{
		Name: "Cycle Counting",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.AddCycles(4)
			cpu.AddCycles(11)
			assert.Equal(t, "Total cycles should be 15", cpu.Cycles, uint64(15))
		},
	},
	{
		Name: "Halt State Management",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			assert.False(t, "Should not be halted initially", cpu.Halted)
			
			cpu.SetHalt(true)
			assert.True(t, "Should be halted", cpu.Halted)
			
			cpu.SetHalt(false)
			assert.False(t, "Should be resumed", cpu.Halted)
		},
	},
	{
		Name: "Memory and IO Access",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Memory.Write(0x1234, 0x42)
			assert.Equal(t, "Memory read should return written value", cpu.Memory.Read(0x1234), uint8(0x42))
			
			cpu.IO.Out(0x0080, 0xAA)
			assert.Equal(t, "IO read should return written value", cpu.IO.In(0x0080), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: NOP",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// MainTable[0x00] is NOP
			instr := MainTable[0x00]
			cycles := instr.Execute(cpu)
			
			assert.Equal(t, "NOP should take 4 cycles", cycles, 4)
			assert.Equal(t, "Mnemonic should be NOP", instr.Mnemonic, "NOP")
		},
	},
}
