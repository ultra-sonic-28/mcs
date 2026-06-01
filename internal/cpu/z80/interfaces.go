// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

// Memory defines the interface for Z80 memory access.
// The Z80 has a 16-bit address bus (64KB).
type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

// IO defines the interface for Z80 I/O port access.
// Although Z80 has an 8-bit port address in some instructions,
// it uses a 16-bit address bus for I/O (often B is the high byte).
type IO interface {
	In(port uint16) uint8
	Out(port uint16, val uint8)
}
