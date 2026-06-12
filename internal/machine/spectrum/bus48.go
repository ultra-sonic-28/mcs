package spectrum

import (
	"log/slog"
)

// Bus48 implements the memory and I/O bus for the ZX Spectrum 48K.
type Bus48 struct {
	BaseBus
	Memory *Memory48
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
		Memory: NewMemory48(),
	}

	return b
}

// Read returns the byte at the specified memory address.
func (b *Bus48) Read(addr uint16) uint8 {
	return b.Memory.Read(addr)
}

// Read16 returns the 16-bit word at the specified memory address (little-endian).
func (b *Bus48) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr)) | (uint16(b.Read(addr+1)) << 8)
}

// Write stores a byte at the specified memory address.
func (b *Bus48) Write(addr uint16, val uint8) {
	b.Memory.Write(addr, val)
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
	return b.Memory.GetDisplayMemory()
}

func (b *Bus48) IsRom1Active() bool {
	return b.Memory.IsRom1Active()
}
