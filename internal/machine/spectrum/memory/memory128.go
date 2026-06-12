package memory

import (
	"log/slog"
	spectrumrom "mcs/assets/machines/spectrum"
)

// Memory128 implements the memory for the ZX Spectrum 128K.
type Memory128 struct {
	rom [2][16384]uint8
	ram [8][16384]uint8

	currentRom   uint8
	currentRam   uint8
	pagingLocked bool
	activeScreen uint8 // 5 or 7
}

// NewMemory128 creates a new Spectrum 128K memory with the ROMs loaded.
func NewMemory128() *Memory128 {
	m := &Memory128{
		activeScreen: 5,
	}

	// Load ROMs
	if len(spectrumrom.Rom128_0) != 16384 {
		slog.Warn("Spectrum 128K ROM 1 size is unexpected", "expected", 16384, "actual", len(spectrumrom.Rom128_0))
	}
	if len(spectrumrom.Rom128_1) != 16384 {
		slog.Warn("Spectrum 128K ROM 2 size is unexpected", "expected", 16384, "actual", len(spectrumrom.Rom128_1))
	}
	copy(m.rom[0][:], spectrumrom.Rom128_0)
	copy(m.rom[1][:], spectrumrom.Rom128_1)

	return m
}

// Read returns the byte at the specified memory address.
func (m *Memory128) Read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.rom[m.currentRom][addr]
	case addr < 0x8000:
		return m.ram[5][addr-0x4000]
	case addr < 0xC000:
		return m.ram[2][addr-0x8000]
	default:
		return m.ram[m.currentRam][addr-0xC000]
	}
}

// Write stores a byte at the specified memory address.
func (m *Memory128) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x4000:
		// ROM is read-only
		return
	case addr < 0x8000:
		m.ram[5][addr-0x4000] = val
	case addr < 0xC000:
		m.ram[2][addr-0x8000] = val
	default:
		m.ram[m.currentRam][addr-0xC000] = val
	}
}

// Page handles memory paging via port 0x7FFD.
func (m *Memory128) Page(val uint8) {
	if m.pagingLocked {
		return
	}

	m.currentRam = val & 0x07
	if (val & 0x08) != 0 {
		m.activeScreen = 7
	} else {
		m.activeScreen = 5
	}
	if (val & 0x10) != 0 {
		m.currentRom = 1
	} else {
		m.currentRom = 0
	}
	if (val & 0x20) != 0 {
		m.pagingLocked = true
	}

	slog.Debug("Spectrum 128K Paging", "ram", m.currentRam, "rom", m.currentRom, "screen", m.activeScreen, "locked", m.pagingLocked)
}

// GetDisplayMemory returns the 6912 bytes of memory used for the display.
func (m *Memory128) GetDisplayMemory() []byte {
	return m.ram[m.activeScreen][0:6912]
}

// IsRom1Active returns true if the 48K BASIC ROM is currently active (ROM 1 on 128K).
func (m *Memory128) IsRom1Active() bool {
	return m.currentRom == 1
}
