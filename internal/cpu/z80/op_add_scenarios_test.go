package z80

import (
	"fmt"
	"mcs/testutils/assert"
	"testing"
)

var addScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: ADD A, n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test: 0x10 + 0x20 = 0x30
			cpu.Regs.A = 0x10
			cpu.Memory.Write(0, 0xC6) // ADD A, n
			cpu.Memory.Write(1, 0x20)
			
			opcode := cpu.FetchByte()
			instr := MainTable[opcode]
			cycles := instr.Execute(cpu)
			
			assert.Equal(t, "A should be 0x30", cpu.Regs.A, uint8(0x30))
			assert.Equal(t, "Cycles should be 7", cycles, 7)
			assert.False(t, "Carry should be false", cpu.Regs.Flag(FlagC))
			assert.False(t, "Zero should be false", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: ADD A, B (with Carry)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test: 0xFF + 0x01 = 0x00 (Carry set, Zero set)
			cpu.Regs.A = 0xFF
			cpu.Regs.B = 0x01
			
			instr := MainTable[0x80] // ADD A, B
			cycles := instr.Execute(cpu)
			
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
			assert.Equal(t, "Cycles should be 4", cycles, 4)
			assert.True(t, "Carry should be true", cpu.Regs.Flag(FlagC))
			assert.True(t, "Zero should be true", cpu.Regs.Flag(FlagZ))
			assert.False(t, "Add/Sub should be false", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: ADD A, (HL) (with Half-Carry)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test: 0x0F + 0x01 = 0x10 (Half-Carry set)
			cpu.Regs.A = 0x0F
			cpu.Regs.SetHL(0x2000)
			bus.Write(0x2000, 0x01)
			
			instr := MainTable[0x86] // ADD A, (HL)
			instr.Execute(cpu)
			
			assert.Equal(t, "A should be 0x10", cpu.Regs.A, uint8(0x10))
			assert.True(t, "Half-Carry should be true", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: ADD A, r (Remaining Registers)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			// regsAdd opcodes: 0x81 (C), 0x82 (D), 0x83 (E), 0x84 (H), 0x85 (L), 0x87 (A)
			// (B is already tested in "ADD A, B (with Carry)")
			tests := []struct {
				op   uint8
				name string
				set  func(r *Registers, val uint8)
			}{
				{0x81, "C", func(r *Registers, val uint8) { r.C = val }},
				{0x82, "D", func(r *Registers, val uint8) { r.D = val }},
				{0x83, "E", func(r *Registers, val uint8) { r.E = val }},
				{0x84, "H", func(r *Registers, val uint8) { r.H = val }},
				{0x85, "L", func(r *Registers, val uint8) { r.L = val }},
				{0x87, "A", func(r *Registers, val uint8) { r.A = val }},
			}

			for _, tt := range tests {
				cpu.Regs.A = 0x10
				if tt.name != "A" {
					tt.set(cpu.Regs, 0x20)
				}

				instr := MainTable[tt.op]
				instr.Execute(cpu)

				var expected uint8
				if tt.name == "A" {
					expected = 0x20 // 0x10 + 0x10
				} else {
					expected = 0x30 // 0x10 + 0x20
				}

				assert.Equal(t, fmt.Sprintf("A should be correct after ADD A, %s", tt.name), cpu.Regs.A, expected)
			}
		},
	},
	{
		Name: "Instruction Execution: ADD HL, BC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test: 0x1000 + 0x2000 = 0x3000
			cpu.Regs.SetHL(0x1000)
			cpu.Regs.SetBC(0x2000)
			
			instr := MainTable[0x09] // ADD HL, BC
			cycles := instr.Execute(cpu)
			
			assert.Equal(t, "HL should be 0x3000", cpu.Regs.HL(), uint16(0x3000))
			assert.Equal(t, "Cycles should be 11", cycles, 11)
			assert.False(t, "Carry should be false", cpu.Regs.Flag(FlagC))
			assert.False(t, "Add/Sub should be false", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: ADD HL, DE (with Carry)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test: 0x8000 + 0x8000 = 0x0000 (Carry set)
			cpu.Regs.SetHL(0x8000)
			cpu.Regs.SetDE(0x8000)
			
			instr := MainTable[0x19] // ADD HL, DE
			instr.Execute(cpu)
			
			assert.Equal(t, "HL should be 0x0000", cpu.Regs.HL(), uint16(0x0000))
			assert.True(t, "Carry should be true", cpu.Regs.Flag(FlagC))
		},
	},
	{
		Name: "Instruction Execution: ADD HL, HL (with Half-Carry)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test: 0x0800 + 0x0800 = 0x1000 (Half-Carry set)
			cpu.Regs.SetHL(0x0800)
			
			instr := MainTable[0x29] // ADD HL, HL
			instr.Execute(cpu)
			
			assert.Equal(t, "HL should be 0x1000", cpu.Regs.HL(), uint16(0x1000))
			assert.True(t, "Half-Carry should be true", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: ADD HL, SP",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			cpu.Regs.SetHL(0x1234)
			cpu.Regs.SP = 0x1111
			
			instr := MainTable[0x39] // ADD HL, SP
			instr.Execute(cpu)
			
			assert.Equal(t, "HL should be 0x2345", cpu.Regs.HL(), uint16(0x2345))
		},
	},
	{
		Name: "Instruction Execution: ADD IX, BC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x1111
			cpu.Regs.SetBC(0x2222)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x09)
			
			cpu.Step()
			
			assert.Equal(t, "IX should be 0x3333", cpu.Regs.IX, uint16(0x3333))
		},
	},
	{
		Name: "Instruction Execution: ADD IX, IX",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x1000
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x29)
			
			cpu.Step()
			
			assert.Equal(t, "IX should be 0x2000", cpu.Regs.IX, uint16(0x2000))
		},
	},
	{
		Name: "Instruction Execution: ADD IY, DE",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x4444
			cpu.Regs.SetDE(0x1111)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x19)
			
			cpu.Step()
			
			assert.Equal(t, "IY should be 0x5555", cpu.Regs.IY, uint16(0x5555))
		},
	},
	{
		Name: "Instruction Execution: ADD IY, IY",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x2000
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x29)
			
			cpu.Step()
			
			assert.Equal(t, "IY should be 0x4000", cpu.Regs.IY, uint16(0x4000))
		},
	},
	{
		Name: "Instruction Execution: ADD A, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIXH(0x20)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x84)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x30", cpu.Regs.A, uint8(0x30))
		},
	},
	{
		Name: "Instruction Execution: ADD A, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIXL(0x05)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x85)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x15", cpu.Regs.A, uint8(0x15))
		},
	},
	{
		Name: "Instruction Execution: ADD A, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.IX = 0x2000
			bus.Write(0x2005, 0x05)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x86)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x15", cpu.Regs.A, uint8(0x15))
		},
	},
	{
		Name: "Instruction Execution: ADD A, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIYH(0x20)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x84)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x30", cpu.Regs.A, uint8(0x30))
		},
	},
	{
		Name: "Instruction Execution: ADD A, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetIYL(0x05)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x85)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x15", cpu.Regs.A, uint8(0x15))
		},
	},
	{
		Name: "Instruction Execution: ADD A, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.IY = 0x2000
			bus.Write(0x2005, 0x05)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x86)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x15", cpu.Regs.A, uint8(0x15))
		},
	},
}
