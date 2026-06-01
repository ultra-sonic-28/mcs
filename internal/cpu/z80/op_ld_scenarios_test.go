package z80

import (
	"fmt"
	"mcs/testutils/assert"
	"testing"
)

var ldScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: LD A, n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Setup memory: 0x3E 0x42 (LD A, 0x42)
			cpu.Memory.Write(0, 0x3E)
			cpu.Memory.Write(1, 0x42)
			
			opcode := cpu.FetchByte()
			instr := MainTable[opcode]
			cycles := instr.Execute(cpu)
			
			assert.Equal(t, "Cycles should be 7", cycles, 7)
			assert.Equal(t, "A should be 0x42", cpu.Regs.A, uint8(0x42))
			assert.Equal(t, "PC should be 2", cpu.Regs.PC, uint16(2))
		},
	},
	{
		Name: "Instruction Execution: LD A, (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x1000)
			bus.Write(0x1000, 0x99)
			
			instr := MainTable[0x7E]
			instr.Execute(cpu)
			
			assert.Equal(t, "A should be 0x99", cpu.Regs.A, uint8(0x99))
		},
	},
	{
		Name: "Instruction Execution: LD A, r (Remaining Registers)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			// regs8 opcodes: 0x79 (C), 0x7A (D), 0x7B (E), 0x7C (H), 0x7D (L), 0x7F (A)
			tests := []struct {
				op   uint8
				name string
				set  func(r *Registers, val uint8)
			}{
				{0x79, "C", func(r *Registers, val uint8) { r.C = val }},
				{0x7A, "D", func(r *Registers, val uint8) { r.D = val }},
				{0x7B, "E", func(r *Registers, val uint8) { r.E = val }},
				{0x7C, "H", func(r *Registers, val uint8) { r.H = val }},
				{0x7D, "L", func(r *Registers, val uint8) { r.L = val }},
				{0x7F, "A", func(r *Registers, val uint8) { r.A = val }},
			}

			for _, tt := range tests {
				if tt.name != "A" {
					cpu.Regs.A = 0x00
					tt.set(cpu.Regs, 0x42)
				} else {
					cpu.Regs.A = 0x42
				}

				instr := MainTable[tt.op]
				instr.Execute(cpu)

				assert.Equal(t, fmt.Sprintf("A should be 0x42 after LD A, %s", tt.name), cpu.Regs.A, uint8(0x42))
			}
		},
	},
	{
		Name: "Instruction Execution: LD B, C",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.C = 0x77
			bus.Write(0x0000, 0x41) // LD B, C
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0x77", cpu.Regs.B, uint8(0x77))
		},
	},
	{
		Name: "Instruction Execution: LD (HL), B",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0xC000)
			cpu.Regs.B = 0xEE
			bus.Write(0x0000, 0x70) // LD (HL), B
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0xC000 should be 0xEE", bus.Read(0xC000), uint8(0xEE))
		},
	},
	{
		Name: "Instruction Execution: LD B, (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0xD000)
			bus.Write(0xD000, 0x55)
			bus.Write(0x0000, 0x46) // LD B, (HL)
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0x55", cpu.Regs.B, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD (HL), n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0xE000)
			bus.Write(0x0000, 0x36) // LD (HL), n
			bus.Write(0x0001, 0xAA)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0xE000 should be 0xAA", bus.Read(0xE000), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD BC, nn",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0x01) // LD BC, nn
			bus.Write(0x0001, 0x34) // low
			bus.Write(0x0002, 0x12) // high
			
			cpu.Step()
			
			assert.Equal(t, "BC should be 0x1234", cpu.Regs.BC(), uint16(0x1234))
		},
	},
	{
		Name: "Instruction Execution: LD (BC), A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetBC(0x2000)
			cpu.Regs.A = 0x55
			bus.Write(0x0000, 0x02) // LD (BC), A
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2000 should be 0x55", bus.Read(0x2000), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD A, (BC)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetBC(0x3000)
			bus.Write(0x3000, 0xAA)
			bus.Write(0x0000, 0x0A) // LD A, (BC)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xAA", cpu.Regs.A, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD (DE), A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetDE(0x4000)
			cpu.Regs.A = 0x77
			bus.Write(0x0000, 0x12) // LD (DE), A
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x4000 should be 0x77", bus.Read(0x4000), uint8(0x77))
		},
	},
	{
		Name: "Instruction Execution: LD A, (DE)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetDE(0x5000)
			bus.Write(0x5000, 0xCC)
			bus.Write(0x0000, 0x1A) // LD A, (DE)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xCC", cpu.Regs.A, uint8(0xCC))
		},
	},
	{
		Name: "Instruction Execution: LD (nn), HL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x1234)
			bus.Write(0x0000, 0x22) // LD (nn), HL
			bus.Write(0x0001, 0x00) // low
			bus.Write(0x0002, 0x20) // high (nn = 0x2000)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2000 should be 0x34", bus.Read(0x2000), uint8(0x34))
			assert.Equal(t, "Memory at 0x2001 should be 0x12", bus.Read(0x2001), uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD HL, (nn)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x3000, 0x78)
			bus.Write(0x3001, 0x56)
			bus.Write(0x0000, 0x2A) // LD HL, (nn)
			bus.Write(0x0001, 0x00) // low
			bus.Write(0x0002, 0x30) // high (nn = 0x3000)
			
			cpu.Step()
			
			assert.Equal(t, "HL should be 0x5678", cpu.Regs.HL(), uint16(0x5678))
		},
	},
	{
		Name: "Instruction Execution: LD (nn), A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			bus.Write(0x0000, 0x32) // LD (nn), A
			bus.Write(0x0001, 0x34) // low
			bus.Write(0x0002, 0x12) // high (nn = 0x1234)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x1234 should be 0x55", bus.Read(0x1234), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD A, (nn)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x5678, 0xAA)
			bus.Write(0x0000, 0x3A) // LD A, (nn)
			bus.Write(0x0001, 0x78) // low
			bus.Write(0x0002, 0x56) // high (nn = 0x5678)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xAA", cpu.Regs.A, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD SP, HL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0xABCD)
			cpu.Regs.SP = 0x0000
			bus.Write(0x0000, 0xF9) // LD SP, HL
			
			cpu.Step()
			
			assert.Equal(t, "SP should be 0xABCD", cpu.Regs.SP, uint16(0xABCD))
			},
			},
			{
			Name: "Instruction Execution: LD IX, nn",
			Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x21) // LD IX, nn
			bus.Write(0x0002, 0x34) // low
			bus.Write(0x0003, 0x12) // high

			cpu.Step()

			assert.Equal(t, "IX should be 0x1234", cpu.Regs.IX, uint16(0x1234))
			},
			},
			{
			Name: "Instruction Execution: LD (nn), IX",
			Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0xABCD
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x22) // LD (nn), IX
			bus.Write(0x0002, 0x00) // low
			bus.Write(0x0003, 0x20) // high (nn = 0x2000)

			cpu.Step()

			assert.Equal(t, "Memory at 0x2000 should be 0xCD", bus.Read(0x2000), uint8(0xCD))
			assert.Equal(t, "Memory at 0x2001 should be 0xAB", bus.Read(0x2001), uint8(0xAB))
			},
			},
			{
			Name: "Instruction Execution: LD IX, (nn)",
			Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x3000, 0x78)
			bus.Write(0x3001, 0x56)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x2A) // LD IX, (nn)
			bus.Write(0x0002, 0x00) // low
			bus.Write(0x0003, 0x30) // high (nn = 0x3000)

			cpu.Step()

			assert.Equal(t, "IX should be 0x5678", cpu.Regs.IX, uint16(0x5678))
			},
			},
			{
			Name: "Instruction Execution: LD IXH, n",
			Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x00FF
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x26) // LD IXH, n
			bus.Write(0x0002, 0x12) // n = 0x12

			cpu.Step()

			assert.Equal(t, "IX should be 0x12FF", cpu.Regs.IX, uint16(0x12FF))
			},
			},
			{
			Name: "Instruction Execution: LD IXL, n",
			Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0xFF00
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x2E) // LD IXL, n
			bus.Write(0x0002, 0x34) // n = 0x34

			cpu.Step()

			assert.Equal(t, "IX should be 0xFF34", cpu.Regs.IX, uint16(0xFF34))
			},
			},
			{
			Name: "Instruction Execution: LD (IX+d), n",
			Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x4000
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x36) // LD (IX+d), n
			bus.Write(0x0002, 0x05) // d = 5
			bus.Write(0x0003, 0x55) // n = 0x55

			cpu.Step()

			assert.Equal(t, "Memory at 0x4005 should be 0x55", bus.Read(0x4005), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD B, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXH(0x12)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x44) // LD B, IXH
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0x12", cpu.Regs.B, uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD B, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXL(0x34)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x45) // LD B, IXL
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0x34", cpu.Regs.B, uint8(0x34))
		},
	},
	{
		Name: "Instruction Execution: LD B, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x5000
			bus.Write(0x5005, 0xAA)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x46) // LD B, (IX+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0xAA", cpu.Regs.B, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD C, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXH(0x56)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x4C) // LD C, IXH
			
			cpu.Step()
			
			assert.Equal(t, "C should be 0x56", cpu.Regs.C, uint8(0x56))
		},
	},
	{
		Name: "Instruction Execution: LD C, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXL(0x78)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x4D) // LD C, IXL
			
			cpu.Step()
			
			assert.Equal(t, "C should be 0x78", cpu.Regs.C, uint8(0x78))
		},
	},
	{
		Name: "Instruction Execution: LD C, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x6000
			bus.Write(0x600A, 0xBB)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x4E) // LD C, (IX+d)
			bus.Write(0x0002, 0x0A) // d = 10
			
			cpu.Step()
			
			assert.Equal(t, "C should be 0xBB", cpu.Regs.C, uint8(0xBB))
		},
	},
	{
		Name: "Instruction Execution: LD D, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXH(0x11)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x54) // LD D, IXH
			
			cpu.Step()
			
			assert.Equal(t, "D should be 0x11", cpu.Regs.D, uint8(0x11))
		},
	},
	{
		Name: "Instruction Execution: LD D, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXL(0x22)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x55) // LD D, IXL
			
			cpu.Step()
			
			assert.Equal(t, "D should be 0x22", cpu.Regs.D, uint8(0x22))
		},
	},
	{
		Name: "Instruction Execution: LD D, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x7000
			bus.Write(0x7002, 0xCC)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x56) // LD D, (IX+d)
			bus.Write(0x0002, 0x02) // d = 2
			
			cpu.Step()
			
			assert.Equal(t, "D should be 0xCC", cpu.Regs.D, uint8(0xCC))
		},
	},
	{
		Name: "Instruction Execution: LD E, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXH(0x33)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x5C) // LD E, IXH
			
			cpu.Step()
			
			assert.Equal(t, "E should be 0x33", cpu.Regs.E, uint8(0x33))
		},
	},
	{
		Name: "Instruction Execution: LD E, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXL(0x44)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x5D) // LD E, IXL
			
			cpu.Step()
			
			assert.Equal(t, "E should be 0x44", cpu.Regs.E, uint8(0x44))
		},
	},
	{
		Name: "Instruction Execution: LD E, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x8000
			bus.Write(0x800F, 0xDD)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x5E) // LD E, (IX+d)
			bus.Write(0x0002, 0x0F) // d = 15
			
			cpu.Step()
			
			assert.Equal(t, "E should be 0xDD", cpu.Regs.E, uint8(0xDD))
		},
	},
	{
		Name: "Instruction Execution: LD IXH, r (Register Components)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test LD IXH, B
			cpu.Regs.B = 0x55
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x60)
			cpu.Step()
			assert.Equal(t, "IXH should be 0x55", cpu.Regs.IXH(), uint8(0x55))

			// Test LD IXH, IXL
			cpu.Regs.SetIXL(0xAA)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xDD)
			bus.Write(0x1001, 0x65)
			cpu.Step()
			assert.Equal(t, "IXH should be 0xAA", cpu.Regs.IXH(), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD IXL, r (Register Components)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			
			// Test LD IXL, C
			cpu.Regs.C = 0x12
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x69)
			cpu.Step()
			assert.Equal(t, "IXL should be 0x12", cpu.Regs.IXL(), uint8(0x12))

			// Test LD IXL, IXH
			cpu.Regs.SetIXH(0x34)
			cpu.Regs.PC = 0x1000
			bus.Write(0x1000, 0xDD)
			bus.Write(0x1001, 0x6C)
			cpu.Step()
			assert.Equal(t, "IXL should be 0x34", cpu.Regs.IXL(), uint8(0x34))
		},
	},
	{
		Name: "Instruction Execution: LD IXL, A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x6F) // LD IXL, A
			
			cpu.Step()
			
			assert.Equal(t, "IXL should be 0xAA", cpu.Regs.IXL(), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD L, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x1000
			bus.Write(0x1005, 0x55)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x6E) // LD L, (IX+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "L should be 0x55", cpu.Regs.L, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD (IX+d), r",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x2000
			cpu.Regs.B = 0x99
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x70) // LD (IX+d), B
			bus.Write(0x0002, 0x0A) // d = 10
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x200A should be 0x99", bus.Read(0x200A), uint8(0x99))
		},
	},
	{
		Name: "Instruction Execution: LD A, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXH(0x55)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x7C) // LD A, IXH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD A, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIXL(0xAA)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x7D) // LD A, IXL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xAA", cpu.Regs.A, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD A, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x3000
			bus.Write(0x3005, 0x12)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x7E) // LD A, (IX+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x12", cpu.Regs.A, uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD SP, IX",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x1234
			cpu.Regs.SP = 0x0000
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xF9) // LD SP, IX
			
			cpu.Step()
			
			assert.Equal(t, "SP should be 0x1234", cpu.Regs.SP, uint16(0x1234))
		},
	},
	{
		Name: "Instruction Execution: LD SP, IY",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x5678
			cpu.Regs.SP = 0x0000
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xF9) // LD SP, IY
			
			cpu.Step()
			
			assert.Equal(t, "SP should be 0x5678", cpu.Regs.SP, uint16(0x5678))
		},
	},
	{
		Name: "Instruction Execution: LD IY, nn",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x21) // LD IY, nn
			bus.Write(0x0002, 0x34) // low
			bus.Write(0x0003, 0x12) // high

			cpu.Step()

			assert.Equal(t, "IY should be 0x1234", cpu.Regs.IY, uint16(0x1234))
		},
	},
	{
		Name: "Instruction Execution: LD (nn), IY",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0xABCD
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x22) // LD (nn), IY
			bus.Write(0x0002, 0x00) // low
			bus.Write(0x0003, 0x20) // high (nn = 0x2000)

			cpu.Step()

			assert.Equal(t, "Memory at 0x2000 should be 0xCD", bus.Read(0x2000), uint8(0xCD))
			assert.Equal(t, "Memory at 0x2001 should be 0xAB", bus.Read(0x2001), uint8(0xAB))
		},
	},
	{
		Name: "Instruction Execution: LD IY, (nn)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x3000, 0x78)
			bus.Write(0x3001, 0x56)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x2A) // LD IY, (nn)
			bus.Write(0x0002, 0x00) // low
			bus.Write(0x0003, 0x30) // high (nn = 0x3000)

			cpu.Step()

			assert.Equal(t, "IY should be 0x5678", cpu.Regs.IY, uint16(0x5678))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x00FF
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x26) // LD IYH, n
			bus.Write(0x0002, 0x12) // n = 0x12

			cpu.Step()

			assert.Equal(t, "IY should be 0x12FF", cpu.Regs.IY, uint16(0x12FF))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0xFF00
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x2E) // LD IYL, n
			bus.Write(0x0002, 0x34) // n = 0x34

			cpu.Step()

			assert.Equal(t, "IY should be 0xFF34", cpu.Regs.IY, uint16(0xFF34))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x4000
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x36) // LD (IY+d), n
			bus.Write(0x0002, 0x05) // d = 5
			bus.Write(0x0003, 0x55) // n = 0x55

			cpu.Step()

			assert.Equal(t, "Memory at 0x4005 should be 0x55", bus.Read(0x4005), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD B, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYH(0x12)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x44) // LD B, IYH
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0x12", cpu.Regs.B, uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD B, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYL(0x34)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x45) // LD B, IYL
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0x34", cpu.Regs.B, uint8(0x34))
		},
	},
	{
		Name: "Instruction Execution: LD B, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x5000
			bus.Write(0x5005, 0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x46) // LD B, (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "B should be 0xAA", cpu.Regs.B, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, B",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.B = 0x55
			cpu.Regs.IY = 0x00FF
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x60) // LD IYH, B
			
			cpu.Step()
			
			assert.Equal(t, "IYH should be 0x55", cpu.Regs.IYH(), uint8(0x55))
			assert.Equal(t, "IY should be 0x55FF", cpu.Regs.IY, uint16(0x55FF))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, B",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.B = 0xAA
			cpu.Regs.IY = 0xFF00
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x68) // LD IYL, B
			
			cpu.Step()
			
			assert.Equal(t, "IYL should be 0xAA", cpu.Regs.IYL(), uint8(0xAA))
			assert.Equal(t, "IY should be 0xFFAA", cpu.Regs.IY, uint16(0xFFAA))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), B",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x3000
			cpu.Regs.B = 0x11
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x70) // LD (IY+d), B
			bus.Write(0x0002, 0x0A) // d = 10
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x300A should be 0x11", bus.Read(0x300A), uint8(0x11))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x55
			cpu.Regs.IY = 0x00FF
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x67) // LD IYH, A
			
			cpu.Step()
			
			assert.Equal(t, "IYH should be 0x55", cpu.Regs.IYH(), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0xAA
			cpu.Regs.IY = 0xFF00
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x6F) // LD IYL, A
			
			cpu.Step()
			
			assert.Equal(t, "IYL should be 0xAA", cpu.Regs.IYL(), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), A",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x2000
			cpu.Regs.A = 0x99
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x77) // LD (IY+d), A
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2005 should be 0x99", bus.Read(0x2005), uint8(0x99))
		},
	},
	{
		Name: "Instruction Execution: LD A, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYH(0x55)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x7C) // LD A, IYH
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD A, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYL(0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x7D) // LD A, IYL
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0xAA", cpu.Regs.A, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD A, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x3000
			bus.Write(0x3005, 0x12)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x7E) // LD A, (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x12", cpu.Regs.A, uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD C, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYH(0x12)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x4C) // LD C, IYH
			
			cpu.Step()
			
			assert.Equal(t, "C should be 0x12", cpu.Regs.C, uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD C, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYL(0x34)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x4D) // LD C, IYL
			
			cpu.Step()
			
			assert.Equal(t, "C should be 0x34", cpu.Regs.C, uint8(0x34))
		},
	},
	{
		Name: "Instruction Execution: LD C, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x5000
			bus.Write(0x5005, 0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x4E) // LD C, (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "C should be 0xAA", cpu.Regs.C, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, C",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.C = 0x55
			cpu.Regs.IY = 0x00FF
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x61) // LD IYH, C
			
			cpu.Step()
			
			assert.Equal(t, "IYH should be 0x55", cpu.Regs.IYH(), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, C",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.C = 0xAA
			cpu.Regs.IY = 0xFF00
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x69) // LD IYL, C
			
			cpu.Step()
			
			assert.Equal(t, "IYL should be 0xAA", cpu.Regs.IYL(), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), C",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x2000
			cpu.Regs.C = 0x99
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x71) // LD (IY+d), C
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2005 should be 0x99", bus.Read(0x2005), uint8(0x99))
		},
	},
	{
		Name: "Instruction Execution: LD D, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYH(0x12)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x54) // LD D, IYH
			
			cpu.Step()
			
			assert.Equal(t, "D should be 0x12", cpu.Regs.D, uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD D, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYL(0x34)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x55) // LD D, IYL
			
			cpu.Step()
			
			assert.Equal(t, "D should be 0x34", cpu.Regs.D, uint8(0x34))
		},
	},
	{
		Name: "Instruction Execution: LD D, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x5000
			bus.Write(0x5005, 0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x56) // LD D, (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "D should be 0xAA", cpu.Regs.D, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, D",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.D = 0x55
			cpu.Regs.IY = 0x00FF
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x62) // LD IYH, D
			
			cpu.Step()
			
			assert.Equal(t, "IYH should be 0x55", cpu.Regs.IYH(), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, D",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.D = 0xAA
			cpu.Regs.IY = 0xFF00
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x6A) // LD IYL, D
			
			cpu.Step()
			
			assert.Equal(t, "IYL should be 0xAA", cpu.Regs.IYL(), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), D",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x2000
			cpu.Regs.D = 0x99
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x72) // LD (IY+d), D
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2005 should be 0x99", bus.Read(0x2005), uint8(0x99))
		},
	},
	{
		Name: "Instruction Execution: LD E, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYH(0x12)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x5C) // LD E, IYH
			
			cpu.Step()
			
			assert.Equal(t, "E should be 0x12", cpu.Regs.E, uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD E, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYL(0x34)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x5D) // LD E, IYL
			
			cpu.Step()
			
			assert.Equal(t, "E should be 0x34", cpu.Regs.E, uint8(0x34))
		},
	},
	{
		Name: "Instruction Execution: LD E, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x5000
			bus.Write(0x5005, 0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x5E) // LD E, (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "E should be 0xAA", cpu.Regs.E, uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, E",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.E = 0x55
			cpu.Regs.IY = 0x00FF
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x63) // LD IYH, E
			
			cpu.Step()
			
			assert.Equal(t, "IYH should be 0x55", cpu.Regs.IYH(), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, E",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.E = 0xAA
			cpu.Regs.IY = 0xFF00
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x6B) // LD IYL, E
			
			cpu.Step()
			
			assert.Equal(t, "IYL should be 0xAA", cpu.Regs.IYL(), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), E",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x2000
			cpu.Regs.E = 0x99
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x73) // LD (IY+d), E
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2005 should be 0x99", bus.Read(0x2005), uint8(0x99))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYH(0x55)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x64) // LD IYH, IYH
			
			cpu.Step()
			
			assert.Equal(t, "IYH should still be 0x55", cpu.Regs.IYH(), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD IYH, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYL(0xAA)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x65) // LD IYH, IYL
			
			cpu.Step()
			
			assert.Equal(t, "IYH should be 0xAA", cpu.Regs.IYH(), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD H, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x6000
			bus.Write(0x600A, 0x77)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x66) // LD H, (IY+d)
			bus.Write(0x0002, 0x0A) // d = 10
			
			cpu.Step()
			
			assert.Equal(t, "H should be 0x77", cpu.Regs.H, uint8(0x77))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYH(0x12)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x6C) // LD IYL, IYH
			
			cpu.Step()
			
			assert.Equal(t, "IYL should be 0x12", cpu.Regs.IYL(), uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD IYL, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetIYL(0x34)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x6D) // LD IYL, IYL
			
			cpu.Step()
			
			assert.Equal(t, "IYL should still be 0x34", cpu.Regs.IYL(), uint8(0x34))
		},
	},
	{
		Name: "Instruction Execution: LD L, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x7000
			bus.Write(0x7005, 0x88)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x6E) // LD L, (IY+d)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "L should be 0x88", cpu.Regs.L, uint8(0x88))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), H",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x2000
			cpu.Regs.H = 0xAA
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x74) // LD (IY+d), H
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2005 should be 0xAA", bus.Read(0x2005), uint8(0xAA))
		},
	},
	{
		Name: "Instruction Execution: LD (IY+d), L",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x3000
			cpu.Regs.L = 0x55
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x75) // LD (IY+d), L
			bus.Write(0x0002, 0x0A) // d = 10
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x300A should be 0x55", bus.Read(0x300A), uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: LD (nn), BC",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetBC(0x1234)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x43)
			bus.Write(0x0002, 0x00)
			bus.Write(0x0003, 0x20) // nn = 0x2000
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2000 should be 0x34", bus.Read(0x2000), uint8(0x34))
			assert.Equal(t, "Memory at 0x2001 should be 0x12", bus.Read(0x2001), uint8(0x12))
		},
	},
	{
		Name: "Instruction Execution: LD BC, (nn)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			bus.Write(0x3000, 0x78)
			bus.Write(0x3001, 0x56)
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x4B)
			bus.Write(0x0002, 0x00)
			bus.Write(0x0003, 0x30) // nn = 0x3000
			
			cpu.Step()
			
			assert.Equal(t, "BC should be 0x5678", cpu.Regs.BC(), uint16(0x5678))
		},
	},
	{
		Name: "Instruction Execution: LD A, I",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.I = 0x55
			cpu.IFF2 = true
			bus.Write(0x0000, 0xED)
			bus.Write(0x0001, 0x57)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x55", cpu.Regs.A, uint8(0x55))
			assert.True(t, "PV flag should match IFF2", cpu.Regs.Flag(FlagPV))
		},
	},
}
