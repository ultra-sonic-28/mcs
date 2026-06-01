// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

// AddressingMode defines the way an operand is accessed by an instruction.
type AddressingMode int

const (
	AddrModeNone AddressingMode = iota
	AddrModeImplied             // Operand is inherent to the instruction (e.g., NOP, HALT)
	AddrModeAccumulator         // Operand is register A
	AddrModeRegister            // Operand is an 8-bit register (B, C, D, E, H, L)
	AddrModeRegisterPair        // Operand is a 16-bit register pair (BC, DE, HL, SP)
	AddrModeImmediate           // Operand is an 8-bit immediate value (n)
	AddrModeImmediate16         // Operand is a 16-bit immediate value (nn)
	AddrModeIndirect            // Operand is a memory location pointed to by a register (e.g., (HL))
	AddrModeExtended            // Operand is a memory location at a 16-bit address (nn)
	AddrModeIndexed             // Operand is a memory location at (IX+d) or (IY+d)
	AddrModeRelative            // Operand is a signed 8-bit displacement from PC
	AddrModePort                // Operand is an I/O port
	AddrModeBit                 // Operand is a bit index (0-7)
)

// String returns a string representation of the addressing mode.
func (m AddressingMode) String() string {
	switch m {
	case AddrModeNone:
		return "None"
	case AddrModeImplied:
		return "Implied"
	case AddrModeAccumulator:
		return "Accumulator"
	case AddrModeRegister:
		return "Register"
	case AddrModeRegisterPair:
		return "Register Pair"
	case AddrModeImmediate:
		return "Immediate"
	case AddrModeImmediate16:
		return "Immediate 16"
	case AddrModeIndirect:
		return "Indirect"
	case AddrModeExtended:
		return "Extended"
	case AddrModeIndexed:
		return "Indexed"
	case AddrModeRelative:
		return "Relative"
	case AddrModePort:
		return "Port"
	case AddrModeBit:
		return "Bit"
	default:
		return "Unknown"
	}
}
