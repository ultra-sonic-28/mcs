// Package bus provides common bus implementations for CPU-to-peripheral communication.
package bus

import "log/slog"

// SimpleBus is a basic implementation of a 64KB flat memory and null I/O.
type SimpleBus struct {
	Memory [65536]uint8
}

// NewSimpleBus creates and initializes a new SimpleBus.
func NewSimpleBus() *SimpleBus {
	slog.Debug("Initializing simple 64KB bus")
	return &SimpleBus{}
}

// Read returns the byte at the specified memory address.
func (b *SimpleBus) Read(addr uint16) uint8 {
	return b.Memory[addr]
}

// Write stores a byte at the specified memory address.
func (b *SimpleBus) Write(addr uint16, val uint8) {
	b.Memory[addr] = val
}

// In reads a byte from the specified I/O port (currently returns 0).
func (b *SimpleBus) In(port uint16) uint8 {
	return 0
}

// Out writes a byte to the specified I/O port (currently a no-op).
func (b *SimpleBus) Out(port uint16, val uint8) {
	// No-op for the basic implementation
}
