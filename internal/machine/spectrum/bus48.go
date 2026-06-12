package spectrum

import (
	"log/slog"
	spectrumrom "mcs/assets/machines/spectrum"
)

// Bus48 implements the memory and I/O bus for the ZX Spectrum 48K.
type Bus48 struct {
	BaseBus
	rom [16384]uint8
	ram [49152]uint8
}

// NewBus48 creates a new Spectrum 48K Bus with the embedded ROM loaded.
func NewBus48() *Bus48 {
	slog.Info("Initializing Spectrum 48K Bus")
	b := &Bus48{
		BaseBus: BaseBus{
			Keyboard: NewKeyboard(),
			Display:  NewDisplay(),
			Tape:     NewTape(),
		},
	}

	// Load ROM
	romData := spectrumrom.Rom48
	if len(romData) != 16384 {
		slog.Warn("Spectrum 48K ROM size is unexpected", "expected", 16384, "actual", len(romData))
	}
	copy(b.rom[:], romData)

	return b
}

// Read returns the byte at the specified memory address.
func (b *Bus48) Read(addr uint16) uint8 {
	if addr < 16384 {
		return b.rom[addr]
	}
	return b.ram[addr-16384]
}

// Read16 returns the 16-bit word at the specified memory address (little-endian).
func (b *Bus48) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr)) | (uint16(b.Read(addr+1)) << 8)
}

// Write stores a byte at the specified memory address.
func (b *Bus48) Write(addr uint16, val uint8) {
	if addr < 16384 {
		return
	}
	b.ram[addr-16384] = val
}

// In reads a byte from the specified I/O port.
func (b *Bus48) In(port uint16) uint8 {
	if port&0x01 == 0 {
		mask := uint8(port >> 8)
		result := b.Keyboard.Scan(mask)

		if b.TapeInState {
			result |= 0x40
		} else {
			result &= ^uint8(0x40)
		}

		result |= 0xA0
		return result
	}
	return 0xFF
}

// Out writes a byte to the specified I/O port.
func (b *Bus48) Out(port uint16, val uint8) {
	if port&0x01 == 0 {
		b.BorderColor = val & 0x07
		b.MicState = (val & 0x08) != 0
		b.BeeperState = (val & 0x10) != 0
	}
}

// GetDisplayMemory returns the 6912 bytes of memory used for the display.
func (b *Bus48) GetDisplayMemory() []byte {
	return b.ram[0:6912]
}

func (b *Bus48) IsRom1Active() bool {
	return true
}
