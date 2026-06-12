package spectrum

import (
	"log/slog"
	spectrumrom "mcs/assets/machines/spectrum"
)

// Memory48 implements the memory for the ZX Spectrum 48K.
type Memory48 struct {
	rom [16384]uint8
	ram [49152]uint8
}

// NewMemory48 creates a new Spectrum 48K memory with the ROM loaded.
func NewMemory48() *Memory48 {
	m := &Memory48{}
	
	// Load ROM
	romData := spectrumrom.Rom48
	if len(romData) != 16384 {
		slog.Warn("Spectrum 48K ROM size is unexpected", "expected", 16384, "actual", len(romData))
	}
	copy(m.rom[:], romData)
	
	return m
}

// Read returns the byte at the specified memory address.
func (m *Memory48) Read(addr uint16) uint8 {
	if addr < 16384 {
		return m.rom[addr]
	}
	return m.ram[addr-16384]
}

// Write stores a byte at the specified memory address.
func (m *Memory48) Write(addr uint16, val uint8) {
	if addr < 16384 {
		// ROM is read-only
		return
	}
	m.ram[addr-16384] = val
}

// GetDisplayMemory returns the 6912 bytes of memory used for the display.
func (m *Memory48) GetDisplayMemory() []byte {
	return m.ram[0:6912]
}

// IsRom1Active returns true as 48K only has one ROM.
func (m *Memory48) IsRom1Active() bool {
	return true
}
