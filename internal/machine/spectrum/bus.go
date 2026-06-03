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

	// Keyboard
	Keyboard *Keyboard

	// ULA State
	BorderColor uint8
	BeeperState bool
	MicState    bool
	TapeInState bool

	// Display
	Display *Display

	// Tape
	Tape *Tape
}

// NewBus creates a new Spectrum Bus with the embedded 48K ROM loaded.
func NewBus() *Bus {
	slog.Info("Initializing Spectrum 48K Bus")
	b := &Bus{
		Keyboard: NewKeyboard(),
		Display:  NewDisplay(),
		Tape:     NewTape(),
	}
	
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
// Port 0xFE is the main ULA port for keyboard and tape.
func (b *Bus) In(port uint16) uint8 {
	// Spectrum decoding for port 0xFE: any even address.
	if port&0x01 == 0 {
		mask := uint8(port >> 8)
		result := b.Keyboard.Scan(mask)

		// Bit 6: EAR/Tape input
		if b.TapeInState {
			result |= 0x40
		} else {
			result &= ^uint8(0x40)
		}

		// Bit 5 & 7: Usually 1
		result |= 0xA0

		return result
	}
	return 0xFF
}

// Out writes a byte to the specified I/O port.
// Port 0xFE handles border, beeper and mic.
func (b *Bus) Out(port uint16, val uint8) {
	if port&0x01 == 0 {
		b.BorderColor = val & 0x07
		b.MicState = (val & 0x08) != 0
		b.BeeperState = (val & 0x10) != 0
	}
}

// GetDisplayMemory returns the 6912 bytes of memory used for the display (starting at 0x4000).
func (b *Bus) GetDisplayMemory() []byte {
	return b.ram[0:6912]
}
