// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

// InterruptMode defines the various interrupt modes of the Z80.
type InterruptMode int

const (
	// IM0: Mode 0 - Data on bus is treated as an instruction (usually RST).
	IM0 InterruptMode = iota
	// IM1: Mode 1 - CPU jumps to 0x0038.
	IM1
	// IM2: Mode 2 - Indirect call via a vector table (using I register).
	IM2
)

// String returns a string representation of the interrupt mode.
func (m InterruptMode) String() string {
	switch m {
	case IM0:
		return "IM0"
	case IM1:
		return "IM1"
	case IM2:
		return "IM2"
	default:
		return "Unknown"
	}
}
