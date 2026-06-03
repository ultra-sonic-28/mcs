// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
	"log/slog"
	"mcs/assets/machines/spectrum"
)

// Bus implements the memory and I/O bus for the ZX Spectrum 48K.
type Bus struct {
	rom [16384]uint8
	ram [49152]uint8
}

// NewBus creates a new Spectrum Bus with the embedded 48K ROM loaded.
func NewBus() *Bus {
	slog.Info("Initializing Spectrum 48K Bus")
	b := &Bus{}
	
	// Load ROM
	romData := spectrumrom.Rom48
	if len(romData) != 16384 {
		slog.Warn("Spectrum ROM size is unexpected", "expected", 16384, "actual", len(romData))
	}
	copy(b.rom[:], romData)
	
	return b
}

// Read returns the byte at the specified memory address.
// 0x0000-0x3FFF: ROM (16KB)
// 0x4000-0xFFFF: RAM (48KB)
func (b *Bus) Read(addr uint16) uint8 {
	if addr < 16384 {
		return b.rom[addr]
	}
	return b.ram[addr-16384]
}

// Write stores a byte at the specified memory address.
// ROM is read-only.
func (b *Bus) Write(addr uint16, val uint8) {
	if addr < 16384 {
		// ROM is read-only, ignore writes
		return
	}
	b.ram[addr-16384] = val
}

// In reads a byte from the specified I/O port.
// For now, it returns 0xFF (no peripheral responding).
func (b *Bus) In(port uint16) uint8 {
	return 0xFF
}

// Out writes a byte to the specified I/O port.
// For now, it is a no-op.
func (b *Bus) Out(port uint16, val uint8) {
	// To be implemented: Border color, Beeper, etc.
}
