// Package bus implements the ZX Spectrum bus logic.
package bus

import (
	"log/slog"
	"mcs/internal/machine/spectrum/display"
	"mcs/internal/machine/spectrum/keyboard"
	"mcs/internal/machine/spectrum/memory"
	"mcs/internal/machine/spectrum/tape"
)

// Bus48 implements the memory and I/O bus for the ZX Spectrum 48K.
type Bus48 struct {
	BaseBus
	Memory *memory.Memory48
}

// NewBus48 creates and initializes a new Bus48.
func NewBus48() *Bus48 {
	return &Bus48{
		BaseBus: BaseBus{
			Keyboard: keyboard.NewKeyboard(),
			Display:  display.NewDisplay(),
			Tape:     tape.NewTape(),
		},
		Memory: memory.NewMemory48(),
	}
}

// Read returns the value from the specified memory address.
func (b *Bus48) Read(addr uint16) uint8 {
	return b.Memory.Read(addr)
}

// Write stores the value at the specified memory address.
func (b *Bus48) Write(addr uint16, val uint8) {
	b.Memory.Write(addr, val)
}

// Read16 reads a 16-bit little-endian value from the specified memory address.
func (b *Bus48) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr)) | uint16(b.Read(addr+1))<<8
}

// In handles I/O input from the specified port.
func (b *Bus48) In(port uint16) uint8 {
	// Any even port (A0=0) addresses the ULA.
	if port&0x0001 == 0 {
		// Bit 6: EAR/Tape Input
		// Bits 0-4: Keyboard matrix
		res := b.Keyboard.Scan(uint8(port >> 8))
		if b.TapeInState {
			res |= 0x40
		} else {
			res &= 0xBF
		}
		// Bit 5: Always 1 on 48K/128K? (Actually usually 1)
		res |= 0x20
		// Bit 7: Always 1
		res |= 0x80
		return res
	}

	slog.Debug("Bus48 IN from unknown port", "port", port)
	return 0xFF
}

// Out handles I/O output to the specified port.
func (b *Bus48) Out(port uint16, val uint8) {
	// Any even port (A0=0) addresses the ULA.
	if port&0x0001 == 0 {
		// Bits 0-2: Border color
		b.BorderColor = val & 0x07
		// Bit 3: MIC
		b.MicState = (val & 0x08) != 0
		// Bit 4: Beeper
		b.BeeperState = (val & 0x10) != 0
		return
	}

	slog.Debug("Bus48 OUT to unknown port", "port", port, "val", val)
}

// GetDisplayMemory returns the memory currently being used for display.
func (b *Bus48) GetDisplayMemory() []byte {
	return b.Memory.GetDisplayMemory()
}

// IsRom1Active returns true for 48K as it only has one ROM.
func (b *Bus48) IsRom1Active() bool {
	return b.Memory.IsRom1Active()
}
