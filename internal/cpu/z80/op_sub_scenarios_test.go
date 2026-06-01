package z80

import (
	"mcs/testutils/assert"
	"testing"
)

var subScenarios = []CPUScenario{
	{
		Name: "Instruction Execution: SUB B",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.B = 0x05
			bus.Write(0x0000, 0x90) // SUB B
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x0B", cpu.Regs.A, uint8(0x0B))
			assert.False(t, "Carry should be false", cpu.Regs.Flag(FlagC))
			assert.True(t, "Add/Sub should be true", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: SBC A, n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			cpu.Regs.SetFlag(FlagC, true)
			bus.Write(0x0000, 0xDE) // SBC A, n
			bus.Write(0x0001, 0x05)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x0A", cpu.Regs.A, uint8(0x0A)) // 0x10 - 5 - 1
		},
	},
	{
		Name: "Instruction Execution: CP n",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x10
			bus.Write(0x0000, 0xFE) // CP n
			bus.Write(0x0001, 0x10)
			
			cpu.Step()
			
			assert.Equal(t, "A should remain 0x10", cpu.Regs.A, uint8(0x10))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
			assert.True(t, "Add/Sub should be true", cpu.Regs.Flag(FlagN))
		},
	},
	{
		Name: "Instruction Execution: SUB A, IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIXH(0x10)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x94)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x20", cpu.Regs.A, uint8(0x20))
		},
	},
	{
		Name: "Instruction Execution: SUB A, IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.SetIXL(0x05)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x95)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x10", cpu.Regs.A, uint8(0x10))
		},
	},
	{
		Name: "Instruction Execution: SUB A, (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.IX = 0x2000
			bus.Write(0x2005, 0x05)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0x96)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x10", cpu.Regs.A, uint8(0x10))
		},
	},
	{
		Name: "Instruction Execution: SUB A, IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIYH(0x10)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x94)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x20", cpu.Regs.A, uint8(0x20))
		},
	},
	{
		Name: "Instruction Execution: SUB A, IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.SetIYL(0x05)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x95)
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x10", cpu.Regs.A, uint8(0x10))
		},
	},
	{
		Name: "Instruction Execution: SUB A, (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.IY = 0x2000
			bus.Write(0x2005, 0x05)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0x96)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should be 0x10", cpu.Regs.A, uint8(0x10))
		},
	},
	{
		Name: "Instruction Execution: CP IXH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIXH(0x30)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xBC) // CP IXH
			
			cpu.Step()
			
			assert.Equal(t, "A should remain 0x30", cpu.Regs.A, uint8(0x30))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: CP IXL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIXL(0x10)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xBD) // CP IXL
			
			cpu.Step()
			
			assert.Equal(t, "A should remain 0x30", cpu.Regs.A, uint8(0x30))
			assert.False(t, "Zero flag should be false", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: CP (IX+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.IX = 0x4000
			bus.Write(0x4005, 0x15)
			bus.Write(0x0000, 0xDD)
			bus.Write(0x0001, 0xBE)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should remain 0x15", cpu.Regs.A, uint8(0x15))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: CP IYH",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIYH(0x30)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xBC) // CP IYH
			
			cpu.Step()
			
			assert.Equal(t, "A should remain 0x30", cpu.Regs.A, uint8(0x30))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: CP IYL",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x30
			cpu.Regs.SetIYL(0x10)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xBD) // CP IYL
			
			cpu.Step()
			
			assert.Equal(t, "A should remain 0x30", cpu.Regs.A, uint8(0x30))
			assert.False(t, "Zero flag should be false", cpu.Regs.Flag(FlagZ))
		},
	},
	{
		Name: "Instruction Execution: CP (IY+d)",
		Run: func(t *testing.T) {
			bus := &MockBus{}
			cpu := NewCPU(bus, bus)
			cpu.Regs.A = 0x15
			cpu.Regs.IY = 0x4000
			bus.Write(0x4005, 0x15)
			bus.Write(0x0000, 0xFD)
			bus.Write(0x0001, 0xBE)
			bus.Write(0x0002, 0x05) // d = 5
			
			cpu.Step()
			
			assert.Equal(t, "A should remain 0x15", cpu.Regs.A, uint8(0x15))
			assert.True(t, "Zero flag should be set", cpu.Regs.Flag(FlagZ))
		},
	},
}
