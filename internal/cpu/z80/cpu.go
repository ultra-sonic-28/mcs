// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

import (
	"log/slog"
)

// CPU represents the Z80 processor state and control logic.
type CPU struct {
	// Registers
	Regs *Registers

	// External Buses
	Memory Memory
	IO     IO

	// Interrupts
	IFF1 bool          // Interrupt Flip-Flop 1 (main interrupt enable)
	IFF2 bool          // Interrupt Flip-Flop 2 (storage for IFF1 during NMI)
	NMI  bool          // Non-Maskable Interrupt pending
	INT  bool          // Maskable Interrupt pending
	IM   InterruptMode // Interrupt Mode

	// State Management
	Halted           bool   // CPU is in HALT state
	Cycles           uint64 // Total number of T-cycles executed
	LastDisplacement int8   // Displacement for DDCB/FDCB instructions
}

// NewCPU creates and initializes a new Z80 CPU instance.
func NewCPU(mem Memory, io IO) *CPU {
	slog.Info("Initializing Z80 CPU")
	cpu := &CPU{
		Regs:   NewRegisters(),
		Memory: mem,
		IO:     io,
	}
	cpu.Reset()
	return cpu
}

// Reset performs a hardware reset of the CPU.
// According to Z80 specifications:
// - PC and IR are cleared to 0.
// - SP is set to 0xFFFF.
// - IFF1 and IFF2 are cleared (interrupts disabled).
// - IM is set to 0.
func (c *CPU) Reset() {
	slog.Info("Resetting Z80 CPU")
	c.Regs.PC = 0
	c.Regs.I = 0
	c.Regs.R = 0
	c.Regs.SP = 0xFFFF
	c.IFF1 = false
	c.IFF2 = false
	c.IM = 0
	c.Halted = false
	c.Cycles = 0
	c.NMI = false
	c.INT = false
}

// AddCycles increments the total T-cycle counter.
func (c *CPU) AddCycles(count uint64) {
	c.Cycles += count
}

// FetchByte reads a byte from memory at the current PC and increments PC.
func (c *CPU) FetchByte() uint8 {
	val := c.Memory.Read(c.Regs.PC)
	c.Regs.PC++
	return val
}

// FetchWord reads a 16-bit word (Little-Endian) from memory at the current PC and increments PC by 2.
func (c *CPU) FetchWord() uint16 {
	low := uint16(c.FetchByte())
	high := uint16(c.FetchByte())
	return high<<8 | low
}

// SetHalt changes the halted state of the CPU.
func (c *CPU) SetHalt(halted bool) {
	if c.Halted != halted {
		slog.Debug("CPU halt state changed", "halted", halted)
	}
	c.Halted = halted
}

// Step executes the next instruction pointed by PC.
// It returns the number of T-cycles consumed.
func (c *CPU) Step() int {
	if c.Halted {
		// In halt state, the CPU executes NOP-like cycles until an interrupt occurs.
		c.AddCycles(4)
		return 4
	}

	opcode := c.FetchByte()
	var instr Instruction

	switch opcode {
	case 0xCB:
		opcode = c.FetchByte()
		instr = CBTable[opcode]
	case 0xED:
		opcode = c.FetchByte()
		instr = EDTable[opcode]
	case 0xDD:
		opcode = c.FetchByte()
		if opcode == 0xCB {
			c.LastDisplacement = int8(c.FetchByte())
			opcode = c.FetchByte()
			instr = DDCBTable[opcode]
		} else {
			instr = DDTable[opcode]
		}
	case 0xFD:
		opcode = c.FetchByte()
		if opcode == 0xCB {
			c.LastDisplacement = int8(c.FetchByte())
			opcode = c.FetchByte()
			instr = FDCBTable[opcode]
		} else {
			instr = FDTable[opcode]
		}
	default:
		instr = MainTable[opcode]
	}

	cycles := instr.Execute(c)
	c.AddCycles(uint64(cycles))
	return cycles
}
