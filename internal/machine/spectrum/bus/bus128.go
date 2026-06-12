// Package bus implements the ZX Spectrum bus logic.
package bus

import (
	"log/slog"
	"mcs/internal/machine/spectrum/display"
	"mcs/internal/machine/spectrum/keyboard"
	"mcs/internal/machine/spectrum/memory"
	"mcs/internal/machine/spectrum/sound"
	"mcs/internal/machine/spectrum/tape"
)

// Bus128 implements the memory and I/O bus for the ZX Spectrum 128K.
type Bus128 struct {
	BaseBus
	Memory *memory.Memory128
	Sound  *sound.AY38912
}

// NewBus128 creates and initializes a new Bus128.
func NewBus128() *Bus128 {
	return &Bus128{
		BaseBus: BaseBus{
			Keyboard: keyboard.NewKeyboard(),
			Display:  display.NewDisplay(),
			Tape:     tape.NewTape(),
		},
		Memory: memory.NewMemory128(),
		Sound:  sound.NewAY38912(),
	}
}

// Read returns the value from the specified memory address.
func (b *Bus128) Read(addr uint16) uint8 {
	return b.Memory.Read(addr)
}

// Write stores the value at the specified memory address.
func (b *Bus128) Write(addr uint16, val uint8) {
	b.Memory.Write(addr, val)
}

// Read16 reads a 16-bit little-endian value from the specified memory address.
func (b *Bus128) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr)) | uint16(b.Read(addr+1))<<8
}

// In handles I/O input from the specified port.
func (b *Bus128) In(port uint16) uint8 {
	// AY-3-8912 Read (Port 0xFFFD)
	// A15=1, A14=1, A1=0
	if (port & 0xC002) == 0xC000 {
		return b.Sound.ReadData()
	}

	// Any even port (A0=0) addresses the ULA.
	if port&0x0001 == 0 {
		res := b.Keyboard.Scan(uint8(port >> 8))
		if b.TapeInState {
			res |= 0x40
		} else {
			res &= 0xBF
		}
		res |= 0x20
		res |= 0x80
		return res
	}

	slog.Debug("Bus128 IN from unknown port", "port", port)
	return 0xFF
}

// Out handles I/O output to the specified port.
func (b *Bus128) Out(port uint16, val uint8) {
	// Memory Paging (Port 0x7FFD)
	// A15=0, A1=0
	if (port & 0x8002) == 0 {
		b.Memory.Page(val)
		return
	}

	// AY-3-8912 Address Select (Port 0xFFFD)
	// A15=1, A14=1, A1=0
	if (port & 0xC002) == 0xC000 {
		b.Sound.WriteAddress(val)
		return
	}

	// AY-3-8912 Data Write (Port 0xBFFD)
	// A15=1, A14=0, A1=0
	if (port & 0xC002) == 0x8000 {
		b.Sound.WriteData(val)
		return
	}

	// Any even port (A0=0) addresses the ULA.
	if port&0x0001 == 0 {
		b.BorderColor = val & 0x07
		b.MicState = (val & 0x08) != 0
		b.BeeperState = (val & 0x10) != 0
		return
	}

	slog.Debug("Bus128 OUT to unknown port", "port", port, "val", val)
}

// GetDisplayMemory returns the memory currently being used for display.
func (b *Bus128) GetDisplayMemory() []byte {
	return b.Memory.GetDisplayMemory()
}

// IsRom1Active returns true if the 48K BASIC ROM is currently active (ROM 1 on 128K).
func (b *Bus128) IsRom1Active() bool {
	return b.Memory.IsRom1Active()
}
