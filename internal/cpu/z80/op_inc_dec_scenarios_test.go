package z80

import (
	"fmt"
	"mcs/testutils/assert"
	"testing"
)

var incDecScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: INC r (All 8-bit registers)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			tests := []struct {
				op   uint8
				name string
				set  func(r *Registers, val uint8)
				get  func(r *Registers) uint8
			}{
				{0x04, "B", func(r *Registers, val uint8) { r.B = val }, func(r *Registers) uint8 { return r.B }},
				{0x0C, "C", func(r *Registers, val uint8) { r.C = val }, func(r *Registers) uint8 { return r.C }},
				{0x14, "D", func(r *Registers, val uint8) { r.D = val }, func(r *Registers) uint8 { return r.D }},
				{0x1C, "E", func(r *Registers, val uint8) { r.E = val }, func(r *Registers) uint8 { return r.E }},
				{0x24, "H", func(r *Registers, val uint8) { r.H = val }, func(r *Registers) uint8 { return r.H }},
				{0x2C, "L", func(r *Registers, val uint8) { r.L = val }, func(r *Registers) uint8 { return r.L }},
				{0x3C, "A", func(r *Registers, val uint8) { r.A = val }, func(r *Registers) uint8 { return r.A }},
			}

			for _, tt := range tests {
				cpu.Reset()
				tt.set(cpu.Regs, 0x10)
				cpu.Regs.SetFlag(FlagC, true) // Carry should remain set

				instr := MainTable[tt.op]
				cycles := instr.Execute(cpu)

				assert.Equal(t, fmt.Sprintf("%s should be 0x11", tt.name), tt.get(cpu.Regs), uint8(0x11))
				assert.Equal(t, "Cycles should be 4", cycles, 4)
				assert.True(t, "Carry flag should remain unaffected", cpu.Regs.Flag(FlagC))
				assert.False(t, "Add/Sub flag (N) should be cleared", cpu.Regs.Flag(FlagN))
			}
		},
	},
	{
		Name: "Instruction Execution: INC r (Flags - Zero, Sign, Half-Carry, Overflow)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			// Zero Flag
			cpu.Regs.A = 0xFF
			MainTable[0x3C].Execute(cpu) // INC A
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
			assert.False(t, "Sign flag should be cleared", cpu.Regs.Flag(FlagS))
			assert.True(t, "Half-Carry flag should be set (0x0F -> 0x10 equivalent)", cpu.Regs.Flag(FlagH))
			assert.False(t, "Overflow flag should be cleared", cpu.Regs.Flag(FlagPV))

			// Sign Flag
			cpu.Regs.A = 0x7F
			MainTable[0x3C].Execute(cpu) // INC A
			assert.Equal(t, "A should be 0x80", cpu.Regs.A, uint8(0x80))
			assert.False(t, "Zero flag should be cleared", cpu.Regs.Flag(FlagZ))
			assert.True(t, "Sign flag should be set", cpu.Regs.Flag(FlagS))
			assert.True(t, "Half-Carry flag should be set", cpu.Regs.Flag(FlagH))
			assert.True(t, "Overflow flag should be set (0x7F -> 0x80)", cpu.Regs.Flag(FlagPV))

			// Half-Carry Flag
			cpu.Regs.A = 0x0E
			MainTable[0x3C].Execute(cpu) // INC A
			assert.False(t, "Half-Carry flag should be cleared", cpu.Regs.Flag(FlagH))
			
			cpu.Regs.A = 0x0F
			MainTable[0x3C].Execute(cpu) // INC A
			assert.True(t, "Half-Carry flag should be set", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: DEC r (All 8-bit registers)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			tests := []struct {
				op   uint8
				name string
				set  func(r *Registers, val uint8)
				get  func(r *Registers) uint8
			}{
				{0x05, "B", func(r *Registers, val uint8) { r.B = val }, func(r *Registers) uint8 { return r.B }},
				{0x0D, "C", func(r *Registers, val uint8) { r.C = val }, func(r *Registers) uint8 { return r.C }},
				{0x15, "D", func(r *Registers, val uint8) { r.D = val }, func(r *Registers) uint8 { return r.D }},
				{0x1D, "E", func(r *Registers, val uint8) { r.E = val }, func(r *Registers) uint8 { return r.E }},
				{0x25, "H", func(r *Registers, val uint8) { r.H = val }, func(r *Registers) uint8 { return r.H }},
				{0x2D, "L", func(r *Registers, val uint8) { r.L = val }, func(r *Registers) uint8 { return r.L }},
				{0x3D, "A", func(r *Registers, val uint8) { r.A = val }, func(r *Registers) uint8 { return r.A }},
			}

			for _, tt := range tests {
				cpu.Reset()
				tt.set(cpu.Regs, 0x10)
				cpu.Regs.SetFlag(FlagC, true) // Carry should remain set

				instr := MainTable[tt.op]
				cycles := instr.Execute(cpu)

				assert.Equal(t, fmt.Sprintf("%s should be 0x0F", tt.name), tt.get(cpu.Regs), uint8(0x0F))
				assert.Equal(t, "Cycles should be 4", cycles, 4)
				assert.True(t, "Carry flag should remain unaffected", cpu.Regs.Flag(FlagC))
				assert.True(t, "Add/Sub flag (N) should be set", cpu.Regs.Flag(FlagN))
			}
		},
	},
	{
		Name: "Instruction Execution: DEC r (Flags - Zero, Sign, Half-Carry, Overflow)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			// Zero Flag
			cpu.Regs.A = 0x01
			MainTable[0x3D].Execute(cpu) // DEC A
			assert.Equal(t, "A should be 0x00", cpu.Regs.A, uint8(0x00))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
			assert.False(t, "Sign flag should be cleared", cpu.Regs.Flag(FlagS))
			assert.False(t, "Half-Carry flag should be cleared", cpu.Regs.Flag(FlagH))
			assert.False(t, "Overflow flag should be cleared", cpu.Regs.Flag(FlagPV))

			// Sign Flag
			cpu.Regs.A = 0x00
			MainTable[0x3D].Execute(cpu) // DEC A
			assert.Equal(t, "A should be 0xFF", cpu.Regs.A, uint8(0xFF))
			assert.False(t, "Zero flag should be cleared", cpu.Regs.Flag(FlagZ))
			assert.True(t, "Sign flag should be set", cpu.Regs.Flag(FlagS))
			assert.True(t, "Half-Carry flag should be set (borrow from bit 4)", cpu.Regs.Flag(FlagH))
			assert.False(t, "Overflow flag should be cleared", cpu.Regs.Flag(FlagPV))

			// Half-Carry Flag
			cpu.Regs.A = 0x11
			MainTable[0x3D].Execute(cpu) // DEC A
			assert.False(t, "Half-Carry flag should be cleared", cpu.Regs.Flag(FlagH))
			
			cpu.Regs.A = 0x10
			MainTable[0x3D].Execute(cpu) // DEC A
			assert.True(t, "Half-Carry flag should be set", cpu.Regs.Flag(FlagH))

			// Overflow Flag
			cpu.Regs.A = 0x80
			MainTable[0x3D].Execute(cpu) // DEC A
			assert.Equal(t, "A should be 0x7F", cpu.Regs.A, uint8(0x7F))
			assert.True(t, "Overflow flag should be set (0x80 -> 0x7F)", cpu.Regs.Flag(FlagPV))
			assert.False(t, "Sign flag should be cleared", cpu.Regs.Flag(FlagS))
		},
	},
	{
		Name: "Instruction Execution: INC rr (16-bit registers)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			tests := []struct {
				op   uint8
				name string
				set  func(r *Registers, val uint16)
				get  func(r *Registers) uint16
			}{
				{0x03, "BC", func(r *Registers, val uint16) { r.SetBC(val) }, func(r *Registers) uint16 { return r.BC() }},
				{0x13, "DE", func(r *Registers, val uint16) { r.SetDE(val) }, func(r *Registers) uint16 { return r.DE() }},
				{0x23, "HL", func(r *Registers, val uint16) { r.SetHL(val) }, func(r *Registers) uint16 { return r.HL() }},
				{0x33, "SP", func(r *Registers, val uint16) { r.SP = val }, func(r *Registers) uint16 { return r.SP }},
			}

			for _, tt := range tests {
				cpu.Reset()
				tt.set(cpu.Regs, 0x1234)
				cpu.Regs.F = 0xAA // Set flags to a known value

				instr := MainTable[tt.op]
				cycles := instr.Execute(cpu)

				assert.Equal(t, fmt.Sprintf("%s should be 0x1235", tt.name), tt.get(cpu.Regs), uint16(0x1235))
				assert.Equal(t, "Cycles should be 6", cycles, 6)
				assert.Equal(t, "Flags should remain unaffected", cpu.Regs.F, uint8(0xAA))

				// Test wrap around
				tt.set(cpu.Regs, 0xFFFF)
				instr.Execute(cpu)
				assert.Equal(t, fmt.Sprintf("%s should wrap to 0x0000", tt.name), tt.get(cpu.Regs), uint16(0x0000))
				assert.Equal(t, "Flags should remain unaffected after wrap", cpu.Regs.F, uint8(0xAA))
			}
		},
	},
	{
		Name: "Instruction Execution: DEC rr (16-bit registers)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)

			tests := []struct {
				op   uint8
				name string
				set  func(r *Registers, val uint16)
				get  func(r *Registers) uint16
			}{
				{0x0B, "BC", func(r *Registers, val uint16) { r.SetBC(val) }, func(r *Registers) uint16 { return r.BC() }},
				{0x1B, "DE", func(r *Registers, val uint16) { r.SetDE(val) }, func(r *Registers) uint16 { return r.DE() }},
				{0x2B, "HL", func(r *Registers, val uint16) { r.SetHL(val) }, func(r *Registers) uint16 { return r.HL() }},
				{0x3B, "SP", func(r *Registers, val uint16) { r.SP = val }, func(r *Registers) uint16 { return r.SP }},
			}

			for _, tt := range tests {
				cpu.Reset()
				tt.set(cpu.Regs, 0x1234)
				cpu.Regs.F = 0x55 // Set flags to a known value

				instr := MainTable[tt.op]
				cycles := instr.Execute(cpu)

				assert.Equal(t, fmt.Sprintf("%s should be 0x1233", tt.name), tt.get(cpu.Regs), uint16(0x1233))
				assert.Equal(t, "Cycles should be 6", cycles, 6)
				assert.Equal(t, "Flags should remain unaffected", cpu.Regs.F, uint8(0x55))

				// Test wrap around
				tt.set(cpu.Regs, 0x0000)
				instr.Execute(cpu)
				assert.Equal(t, fmt.Sprintf("%s should wrap to 0xFFFF", tt.name), tt.get(cpu.Regs), uint16(0xFFFF))
				assert.Equal(t, "Flags should remain unaffected after wrap", cpu.Regs.F, uint8(0x55))
			}
		},
	},
	{
		Name: "Instruction Execution: INC (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x2000)
			bus.Write(0x2000, 0x10)
			bus.Write(0x0000, 0x34) // INC (HL)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2000 should be 0x11", bus.Read(0x2000), uint8(0x11))
		},
	},
	{
		Name: "Instruction Execution: DEC (HL)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.SetHL(0x3000)
			bus.Write(0x3000, 0x10)
			bus.Write(0x0000, 0x35) // DEC (HL)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x3000 should be 0x0F", bus.Read(0x3000), uint8(0x0F))
		},
	},
	{
		Name: "Instruction Execution: INC IX",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x1234
			cpu.Regs.F = 0x55
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x23)
			
			cpu.Step()
			
			assert.Equal(t, "IX should be 0x1235", cpu.Regs.IX, uint16(0x1235))
			assert.Equal(t, "Flags should remain unaffected", cpu.Regs.F, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: DEC IX",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x1234
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x2B)
			
			cpu.Step()
			
			assert.Equal(t, "IX should be 0x1233", cpu.Regs.IX, uint16(0x1233))
		},
	},
	{
		Name: "Instruction Execution: INC IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x0F00
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x24)
			
			cpu.Step()
			
			assert.Equal(t, "IXH should be 0x10", cpu.Regs.IXH(), uint8(0x10))
			assert.True(t, "Half-Carry should be set", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: DEC IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x0010
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x2D)
			
			cpu.Step()
			
			assert.Equal(t, "IXL should be 0x0F", cpu.Regs.IXL(), uint8(0x0F))
			assert.True(t, "Half-Carry should be set", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: INC IY",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x1234
			cpu.Regs.F = 0x55
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x23)
			
			cpu.Step()
			
			assert.Equal(t, "IY should be 0x1235", cpu.Regs.IY, uint16(0x1235))
			assert.Equal(t, "Flags should remain unaffected", cpu.Regs.F, uint8(0x55))
		},
	},
	{
		Name: "Instruction Execution: DEC IY",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x1234
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x2B)
			
			cpu.Step()
			
			assert.Equal(t, "IY should be 0x1233", cpu.Regs.IY, uint16(0x1233))
		},
	},
	{
		Name: "Instruction Execution: INC IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x0F00
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x24)
			
			cpu.Step()
			
			assert.Equal(t, "IYH should be 0x10", cpu.Regs.IYH(), uint8(0x10))
			assert.True(t, "Half-Carry should be set", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: DEC IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x0010
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x2D)
			
			cpu.Step()
			
			assert.Equal(t, "IYL should be 0x0F", cpu.Regs.IYL(), uint8(0x0F))
			assert.True(t, "Half-Carry should be set", cpu.Regs.Flag(FlagH))
		},
	},
	{
		Name: "Instruction Execution: INC (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x2000
			bus.Write(0x2005, 0x10)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x34)
			bus.Write(0x0002, 0x05)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2005 should be 0x11", bus.Read(0x2005), uint8(0x11))
			assert.False(t, "Add/Sub flag (N) should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: DEC (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IX = 0x2000
			bus.Write(0x2005, 0x10)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x35)
			bus.Write(0x0002, 0x05)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x2005 should be 0x0F", bus.Read(0x2005), uint8(0x0F))
			assert.True(t, "Add/Sub flag (N) should be set", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: INC (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x3000
			bus.Write(0x300A, 0xFF)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x34)
			bus.Write(0x0002, 0x0A)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x300A should be 0x00", bus.Read(0x300A), uint8(0x00))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
			assert.False(t, "Add/Sub flag (N) should be cleared", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: DEC (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.IY = 0x3000
			bus.Write(0x300A, 0x00)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x35)
			bus.Write(0x0002, 0x0A)
			
			cpu.Step()
			
			assert.Equal(t, "Memory at 0x300A should be 0xFF", bus.Read(0x300A), uint8(0xFF))
			assert.True(t, "Sign flag should be set", cpu.Regs.Flag(FlagS))
			assert.True(t, "Add/Sub flag (N) should be set", cpu.Regs.Flag(FlagN))
		},
	},
}
